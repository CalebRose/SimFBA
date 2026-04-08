# Golang Technical Design: Notifications and Forum Thread/Post Creation with Firebase

## Notes Before Beginning

Please note that before beginning that the files within t_FirebaseLogicFromFrontendCode are based off a design that comes from the inspiration from this design. This document is to help provide an overview on what we need to do with the Golang API through interacting with Firebase. The code within t_FirebaseLogicFromFrontendCode will help provide context as a means of what the code looks like on the Frontend.

## Document Info

- **Project**: Simulation Sports Backend Notifications + Forum Automation
- **Backend Stack**: Golang, Firebase Admin SDK / Firestore
- **Frontend Context**: React.js, TypeScript, Tailwind CSS, Firebase-backed notifications and forums
- **Primary Goals**:
  - Send in-app notifications after sports-domain events
  - Support future email and Discord fan-out without redesigning the core backend
  - Create forum threads and posts automatically from Go services
  - Include explicit routing metadata in each notification so the frontend can navigate users to the correct page

---

## 1. Overview

This document defines a backend-first design for:

- generating notifications from Go services
- storing notifications in Firestore
- creating forum threads and posts from Go APIs
- ensuring notifications include route metadata for frontend navigation
- supporting future delivery channels such as email and Discord

This design assumes:

- the Go API is the source of truth for sports-domain events
- Firebase Authentication manages user identity
- Firestore stores notifications, forum threads, forum posts, and eventually user notification settings
- the frontend already has or will have a notification UI and forum UI backed by Firebase

The design should prioritize:

- idempotency
- explicit routing contracts
- clean separation of domain logic from delivery logic
- easy extension into email and Discord later
- low operational complexity

---

## 2. Goals

### 2.1 Functional Goals

The backend should support notifications for events such as:

- a player is injured
- a recruit signs
- a practice squad player receives an offer
- a recruiting sync completes
- a team needs to update its gameplan

The backend should also support:

- creating system-generated forum threads
- creating system-generated forum posts
- linking notifications to exact frontend destinations
- deep-linking into forum threads and posts
- future fan-out into:
  - email
  - Discord
  - digests or batched notifications

### 2.2 Non-Goals

This version does not require:

- mobile push notifications
- direct SMTP sending from Go
- a Discord listener service
- full user preference management implementation
- full-text search design
- image upload design for forums

---

## 3. Architectural Principles

### 3.1 The Go API Owns Domain Events

The Go backend already knows when meaningful simulation events occur. It should remain responsible for deciding:

- what happened
- who should be notified
- what the notification should say
- where the user should go when they click it
- whether a forum artifact should be created

The frontend should not have to infer this from partial data.

### 3.2 Notifications Should Be Typed and Explicit

Every notification should have:

- a stable type
- a user-facing title
- a user-facing message
- an explicit route object
- metadata describing the referenced entity

Avoid one-off notification payloads that vary unpredictably by service.

### 3.3 Routing Must Be Backend-Defined

Do not rely on the frontend to guess where a notification goes based only on IDs.

Each notification should include:

- `routeName`
- `path`
- `params`
- optional `query`
- optional `anchor`

This keeps routing logic deterministic.

### 3.4 Forum Automation Must Be Idempotent

Automated threads and posts should never duplicate when jobs retry or services rerun.

Examples:

- one postgame thread per game
- one recruiting sync announcement per sync run
- one injury bulletin per source event

Use stable external event keys.

---

## 4. High-Level Architecture

```txt
Sports Domain Services (Go)
    |
    v
Application Event / Domain Event
    |
    v
Notification Service
    |----> Firestore notifications/
    |----> Firestore notificationEvents/ (optional outbox)
    |----> Firestore threads/
    |----> Firestore posts/
    |
    +----> Future emailQueue/
    +----> Future discordQueue/
```

### Backend responsibilities

- detect domain events
- resolve recipients
- generate notification payloads
- create notification docs
- create forum threads/posts when needed
- maintain idempotency

### Frontend responsibilities

- render notifications
- navigate based on notification route metadata
- render forum thread/post destinations
- mark notifications read

---

## 5. Firebase Integration from Go

Recommended server-side options:

- **Firebase Admin SDK for Go**
- **Google Cloud Firestore Go client**

Because this is privileged backend automation, the service should use server credentials and write directly to Firestore.

### Recommendation

Create a shared infrastructure package that initializes:

- Firebase app
- Firestore client
- optional Auth admin client

Suggested package structure:

```txt
internal/
  firebase/
    client.go
  notifications/
    service.go
    repository.go
    builders.go
    types.go
  forums/
    service.go
    repository.go
    types.go
  routing/
    routes.go
  domain/
    injuries/
    recruiting/
    practice/
    gameplan/
```

---

## 6. Firestore Data Model

## 6.1 Notifications Collection

Collection:

```txt
notifications/{notificationId}
```

Suggested document shape:

```go
type Notification struct {
    ID             string                 `firestore:"id"`
    RecipientUID   string                 `firestore:"recipientUid"`
    Type           string                 `firestore:"type"`
    Title          string                 `firestore:"title"`
    Message        string                 `firestore:"message"`
    Severity       string                 `firestore:"severity"`
    IsRead         bool                   `firestore:"isRead"`
    CreatedAt      time.Time              `firestore:"createdAt"`
    ReadAt         *time.Time             `firestore:"readAt,omitempty"`

    League         string                 `firestore:"league,omitempty"`
    TeamID         uint                   `firestore:"teamId,omitempty"`
    SeasonID       uint                   `firestore:"seasonId,omitempty"`
    WeekID         uint                   `firestore:"weekId,omitempty"`

    EntityType     string                 `firestore:"entityType,omitempty"`
    EntityID       string                 `firestore:"entityId,omitempty"`

    Route          NotificationRoute      `firestore:"route"`
    Metadata       map[string]interface{} `firestore:"metadata,omitempty"`
    Delivery       NotificationDelivery   `firestore:"delivery"`

    SourceEventKey string                 `firestore:"sourceEventKey,omitempty"`
}
```

### Why these fields matter

- `Type` lets the frontend render consistent badges/styles
- `Route` lets the frontend navigate without guessing
- `Metadata` supports richer UX without forcing every field into the top-level schema
- `SourceEventKey` supports idempotency and debugging
- `Delivery` leaves room for email and Discord later

---

## 6.2 Notification Route Contract

```go
type NotificationRoute struct {
    RouteName string                 `firestore:"routeName"`
    Path      string                 `firestore:"path"`
    Params    map[string]string      `firestore:"params,omitempty"`
    Query     map[string]string      `firestore:"query,omitempty"`
    Anchor    string                 `firestore:"anchor,omitempty"`
    State     map[string]interface{} `firestore:"state,omitempty"`
}
```

### Example

```go
NotificationRoute{
    RouteName: "team_recruiting",
    Path: "/simcfb/teams/44/recruiting",
    Params: map[string]string{
        "league": "simcfb",
        "teamId": "44",
    },
}
```

The frontend should use `Path` as the primary destination and may use `RouteName`, `Params`, and `State` for enhanced behavior.

---

## 6.3 Notification Delivery State

```go
type NotificationDelivery struct {
    InAppCreated  bool `firestore:"inAppCreated"`
    EmailQueued   bool `firestore:"emailQueued"`
    DiscordQueued bool `firestore:"discordQueued"`
}
```

This allows the same core notification pipeline to expand later without redesigning the schema.

---

## 6.4 Optional Notification Outbox

Collection:

```txt
notificationEvents/{eventId}
```

Suggested document shape:

```go
type NotificationEvent struct {
    ID             string                 `firestore:"id"`
    EventType      string                 `firestore:"eventType"`
    SourceEventKey string                 `firestore:"sourceEventKey"`
    League         string                 `firestore:"league,omitempty"`
    TeamID         uint                   `firestore:"teamId,omitempty"`
    EntityType     string                 `firestore:"entityType,omitempty"`
    EntityID       string                 `firestore:"entityId,omitempty"`
    Payload        map[string]interface{} `firestore:"payload"`
    RecipientUIDs  []string               `firestore:"recipientUids"`
    Status         string                 `firestore:"status"`
    CreatedAt      time.Time              `firestore:"createdAt"`
    ProcessedAt    *time.Time             `firestore:"processedAt,omitempty"`
}
```

### Why keep an outbox?

- easier replay
- better debugging
- easier retries
- simpler fan-out into email or Discord later

### Recommendation

You can write notifications directly for V1. For a cleaner long-term design, add `notificationEvents/`.

---

## 6.5 Threads Collection

Collection:

```txt
threads/{threadId}
```

Suggested document shape:

```go
type ForumThread struct {
    ID               string                 `firestore:"id"`
    ForumID          string                 `firestore:"forumId"`
    ForumPath        []string               `firestore:"forumPath"`
    Title            string                 `firestore:"title"`
    Slug             string                 `firestore:"slug"`
    AuthorUID        string                 `firestore:"authorUid"`
    AuthorName       string                 `firestore:"authorName"`
    CreatedByType    string                 `firestore:"createdByType"` // user, system, bot
    ThreadType       string                 `firestore:"threadType"`    // standard, poll, game_reference, system_event
    FirstPostID      string                 `firestore:"firstPostId"`
    IsPinned         bool                   `firestore:"isPinned"`
    IsLocked         bool                   `firestore:"isLocked"`
    IsDeleted        bool                   `firestore:"isDeleted"`
    ReplyCount       int                    `firestore:"replyCount"`
    ParticipantCount int                    `firestore:"participantCount"`
    LatestPostID     string                 `firestore:"latestPostId,omitempty"`
    LatestActivityAt time.Time              `firestore:"latestActivityAt"`
    CreatedAt        time.Time              `firestore:"createdAt"`
    UpdatedAt        time.Time              `firestore:"updatedAt"`

    ReferencedGameID string                 `firestore:"referencedGameId,omitempty"`
    ExternalEventKey string                 `firestore:"externalEventKey,omitempty"`
    Metadata         map[string]interface{} `firestore:"metadata,omitempty"`
}
```

---

## 6.6 Posts Collection

Collection:

```txt
posts/{postId}
```

Suggested document shape:

```go
type ForumPost struct {
    ID            string                 `firestore:"id"`
    ThreadID      string                 `firestore:"threadId"`
    ForumID       string                 `firestore:"forumId"`
    AuthorUID     string                 `firestore:"authorUid"`
    AuthorName    string                 `firestore:"authorName"`
    CreatedByType string                 `firestore:"createdByType"` // user, system, bot

    Body          map[string]interface{} `firestore:"body"`     // rich text JSON
    BodyText      string                 `firestore:"bodyText"` // plain text fallback
    ReplyToPostID string                 `firestore:"replyToPostId,omitempty"`
    QuotedPostID  string                 `firestore:"quotedPostId,omitempty"`

    IsEdited      bool                   `firestore:"isEdited"`
    IsDeleted     bool                   `firestore:"isDeleted"`

    CreatedAt     time.Time              `firestore:"createdAt"`
    UpdatedAt     time.Time              `firestore:"updatedAt"`

    Metadata      map[string]interface{} `firestore:"metadata,omitempty"`
}
```

---

## 7. Notification Types

Define stable constants in Go.

```go
const (
    NotificationPlayerInjured          = "player_injured"
    NotificationRecruitSigned          = "recruit_signed"
    NotificationPracticeSquadOffer     = "practice_squad_offer"
    NotificationRecruitingSyncComplete = "recruiting_sync_complete"
    NotificationGameplanUpdateNeeded   = "gameplan_update_needed"
    NotificationForumThreadCreated     = "forum_thread_created"
    NotificationForumPostReply         = "forum_post_reply"
)
```

Each type should have:

- a message template
- a title template
- a severity
- a route builder
- a recipient resolution strategy

---

## 8. Notification Route Strategy

The backend should own route generation in one place.

Suggested package:

```txt
internal/routing/routes.go
```

Suggested interface:

```go
type RouteBuilder interface {
    BuildPlayerInjuryRoute(league string, teamID uint, playerID uint) NotificationRoute
    BuildRecruitingRoute(league string, teamID uint) NotificationRoute
    BuildPracticeSquadRoute(league string, teamID uint, playerID uint) NotificationRoute
    BuildGameplanRoute(league string, teamID uint) NotificationRoute
    BuildForumThreadRoute(threadID string) NotificationRoute
    BuildForumPostRoute(threadID string, postID string) NotificationRoute
}
```

### Example implementation

```go
func BuildGameplanRoute(league string, teamID uint) NotificationRoute {
    path := fmt.Sprintf("/%s/teams/%d/gameplan", league, teamID)
    return NotificationRoute{
        RouteName: "team_gameplan",
        Path:      path,
        Params: map[string]string{
            "league": league,
            "teamId": strconv.Itoa(int(teamID)),
        },
    }
}
```

### Recommended mapping

| Notification Type          | Route Name                            | Suggested Path                                   |
| -------------------------- | ------------------------------------- | ------------------------------------------------ |
| `player_injured`           | `team_player_detail` or `team_roster` | `/{league}/teams/{teamId}/roster` or player page |
| `recruit_signed`           | `team_recruiting`                     | `/{league}/teams/{teamId}/recruiting`            |
| `practice_squad_offer`     | `team_practice_squad`                 | `/{league}/teams/{teamId}/practice-squad`        |
| `recruiting_sync_complete` | `recruiting_sync_status`              | sync status/admin page                           |
| `gameplan_update_needed`   | `team_gameplan`                       | `/{league}/teams/{teamId}/gameplan`              |
| `forum_thread_created`     | `forum_thread`                        | `/forums/thread/{threadId}`                      |
| `forum_post_reply`         | `forum_post`                          | `/forums/thread/{threadId}` + anchor             |

Keep `RouteName` stable even if path formats change later.

---

## 9. Recipient Resolution

Do not scatter recipient logic throughout domain services.

Create a resolver layer.

```go
type RecipientResolver interface {
    ResolvePlayerInjuryRecipients(ctx context.Context, teamID uint) ([]string, error)
    ResolveRecruitSignedRecipients(ctx context.Context, teamID uint) ([]string, error)
    ResolvePracticeSquadOfferRecipients(ctx context.Context, teamID uint) ([]string, error)
    ResolveRecruitingSyncRecipients(ctx context.Context, league string) ([]string, error)
    ResolveGameplanRecipients(ctx context.Context, teamID uint) ([]string, error)
}
```

### Expected patterns

- **player injured**
  - team owner
  - coaches
  - team admins
- **recruit signed**
  - team owner
  - recruiting staff
  - team admins
- **practice squad offer**
  - team owner
  - contract manager
  - team admins
- **recruiting sync complete**
  - league ops
  - admins
- **gameplan update needed**
  - team owner
  - active coach users

---

## 10. Notification Service Design

Suggested service interface:

```go
type NotificationService interface {
    NotifyPlayerInjured(ctx context.Context, input PlayerInjuryNotificationInput) error
    NotifyRecruitSigned(ctx context.Context, input RecruitSignedNotificationInput) error
    NotifyPracticeSquadOffer(ctx context.Context, input PracticeSquadOfferNotificationInput) error
    NotifyRecruitingSyncComplete(ctx context.Context, input RecruitingSyncCompleteInput) error
    NotifyGameplanUpdateNeeded(ctx context.Context, input GameplanUpdateNeededInput) error
}
```

### Suggested input types

```go
type PlayerInjuryNotificationInput struct {
    League         string
    TeamID         uint
    TeamName       string
    PlayerID       uint
    PlayerName     string
    InjuryName     string
    GamesMissed    int
    SourceEventKey string
}

type RecruitSignedNotificationInput struct {
    League         string
    TeamID         uint
    TeamName       string
    RecruitID      uint
    RecruitName    string
    SourceEventKey string
}

type PracticeSquadOfferNotificationInput struct {
    League           string
    TeamID           uint
    TeamName         string
    PlayerID         uint
    PlayerName       string
    OfferingTeamID   uint
    OfferingTeamName string
    SourceEventKey   string
}

type RecruitingSyncCompleteInput struct {
    League         string
    SyncRunID      string
    SummaryMessage string
    SourceEventKey string
}

type GameplanUpdateNeededInput struct {
    League         string
    TeamID         uint
    TeamName       string
    Reason         string
    SourceEventKey string
}
```

---

## 11. Notification Build Flow

Each notification method should follow the same shape:

1. validate input
2. check idempotency using `SourceEventKey`
3. resolve recipients
4. build title/message
5. build route
6. create one notification doc per recipient
7. optionally create outbox record

Example flow:

```txt
Game simulation finishes
  -> player injury detected
  -> build PlayerInjuryNotificationInput
  -> notification service checks SourceEventKey
  -> resolve team recipients
  -> build team roster or player route
  -> write notifications
```

---

## 12. Idempotency Strategy

This is critical.

### Why

- job retries happen
- event handlers rerun
- infrastructure can duplicate work
- forum automation is especially sensitive to duplicates

### Requirement

Every domain event that can create notifications or forum artifacts must generate a stable key.

### Example keys

```txt
injury:simphl:season12:game881:player193
recruit_sign:simcfb:season5:team44:recruit1182
practice_offer:simnfl:season4:team12:player892:offerTeam31
recruit_sync:simcfb:sync_2026_04_07T12_00_00Z
gameplan_needed:simchl:season2:team8:week9
postgame_thread:simcfb:season4:game1219
```

### Recommended persistence options

Option A:

- store processed events in `notificationEvents/`

Option B:

- query notifications/threads by `SourceEventKey` or `ExternalEventKey`

### Recommendation

Use an outbox or event registry for clarity.

---

## 13. Forum Automation Design

Suggested service interface:

```go
type ForumService interface {
    CreateThread(ctx context.Context, input CreateForumThreadInput) (*ForumThread, error)
    CreatePost(ctx context.Context, input CreateForumPostInput) (*ForumPost, error)
    FindThreadByExternalEventKey(ctx context.Context, eventKey string) (*ForumThread, error)
}
```

### Create thread input

```go
type CreateForumThreadInput struct {
    ForumID           string
    ForumPath         []string
    Title             string
    Slug              string
    AuthorUID         string
    AuthorName        string
    CreatedByType     string
    ThreadType        string
    FirstPostBody     map[string]interface{}
    FirstPostBodyText string
    ReferencedGameID  string
    ExternalEventKey  string
    Metadata          map[string]interface{}
}
```

### Create post input

```go
type CreateForumPostInput struct {
    ThreadID      string
    ForumID       string
    AuthorUID     string
    AuthorName    string
    CreatedByType string
    Body          map[string]interface{}
    BodyText      string
    ReplyToPostID string
    QuotedPostID  string
    Metadata      map[string]interface{}
}
```

---

## 14. Batched Thread + First Post Creation

When creating a new thread, the backend should create both:

- the thread document
- the first post document

Use a Firestore batch or transaction.

### Recommended flow

1. generate `threadId`
2. generate `postId`
3. create post with `threadId`
4. create thread with `firstPostId = postId`
5. optionally update forum counters/latest activity
6. commit batch

This keeps thread creation atomic enough for the application.

---

## 15. Rich Text Body Format from Go

The frontend may use structured JSON for forum post rendering.

The Go backend should be able to produce:

- a plain text fallback via `BodyText`
- a structured JSON body via `Body`

### Example simple system post

```go
body := map[string]interface{}{
    "type": "doc",
    "content": []map[string]interface{}{
        {
            "type": "paragraph",
            "content": []map[string]interface{}{
                {
                    "type": "text",
                    "text": "Postgame discussion is now open. Share your thoughts here.",
                },
            },
        },
    },
}
```

This allows automated Go-created posts to use the same frontend renderer as user-authored posts.

---

## 16. Example Notification Templates

### 16.1 Player Injured

**Title**  
`Injury Update: {PlayerName}`

**Message**  
`{PlayerName} suffered {InjuryName} and is expected to miss {GamesMissed} games.`

**Route**  
Team roster or player detail page

**Metadata**

```json
{
  "playerId": 193,
  "injuryName": "Separated Shoulder",
  "gamesMissed": 4
}
```

### 16.2 Recruit Signed

**Title**  
`Recruit Signed: {RecruitName}`

**Message**  
`{RecruitName} has signed with {TeamName}.`

**Route**  
Team recruiting page

### 16.3 Practice Squad Offer

**Title**  
`Practice Squad Offer Received`

**Message**  
`{OfferingTeamName} has made an offer for {PlayerName}.`

**Route**  
Practice squad/contracts page

### 16.4 Recruiting Sync Complete

**Title**  
`Recruiting Sync Complete`

**Message**  
`The latest recruiting sync has completed successfully.`

**Route**  
Recruiting admin or sync status page

### 16.5 Gameplan Update Needed

**Title**  
`Gameplan Update Needed`

**Message**  
`Your team needs a gameplan update before the next simulation cycle.`

**Route**  
Team gameplan page

---

## 17. Example Automated Forum Flows

### 17.1 Postgame Thread

**Trigger**  
A game becomes final.

**Flow**

1. build `ExternalEventKey = postgame_thread:{league}:{gameId}`
2. check whether a thread already exists
3. create thread in appropriate forum
4. create first post with matchup summary
5. store referenced game ID
6. optionally notify subscribers or team users

Suggested title:

```txt
Postgame Thread: Away Team at Home Team
```

Suggested thread type:

```txt
game_reference
```

### 17.2 Recruiting Sync Announcement

**Trigger**  
A recruiting sync job finishes.

**Flow**

1. build stable sync event key
2. create thread in admin or recruiting forum, or post into existing operations thread
3. notify admins or league ops users

### 17.3 Injury Bulletin

**Trigger**  
A major injury occurs and should be publicly surfaced.

**Flow**

1. create in-app notifications for team stakeholders
2. optionally create a thread/post in a news forum if the product wants public visibility

---

## 18. Repository Layer

Create repositories so services are not directly coupled to Firestore query details.

### Notification repository

```go
type NotificationRepository interface {
    CreateNotifications(ctx context.Context, notifications []Notification) error
    ExistsBySourceEventKey(ctx context.Context, sourceEventKey string, notificationType string) (bool, error)
    CreateEvent(ctx context.Context, event NotificationEvent) error
    MarkEventProcessed(ctx context.Context, eventID string) error
}
```

### Forum repository

```go
type ForumRepository interface {
    CreateThreadWithFirstPost(ctx context.Context, thread ForumThread, post ForumPost) error
    CreatePost(ctx context.Context, post ForumPost) error
    FindThreadByExternalEventKey(ctx context.Context, eventKey string) (*ForumThread, error)
}
```

### Why use repositories

- easier unit testing
- clear ownership of Firestore logic
- easier future refactors

---

## 19. Frontend Contract Requirements

The frontend notification renderer should expect:

- `type`
- `title`
- `message`
- `route.routeName`
- `route.path`
- `route.params`
- `route.query`
- `route.anchor`
- `metadata`

### Click behavior

1. user clicks notification
2. frontend navigates to `route.path`
3. if `route.anchor` exists, scroll to that element
4. optional state/query can pre-open the right tab or subview

### Example forum reply notification

```json
{
  "type": "forum_post_reply",
  "title": "New reply in Postgame Thread",
  "message": "A new reply was posted in your thread.",
  "route": {
    "routeName": "forum_post",
    "path": "/forums/thread/abc123",
    "anchor": "post_xyz789"
  }
}
```

The frontend should navigate to `/forums/thread/abc123` and scroll to `#post_xyz789`.

---

## 20. Suggested Firestore Collections

```txt
notifications/
notificationEvents/
threads/
posts/
forums/
userNotificationSettings/
emailQueue/
discordQueue/
```

### Notes

- `userNotificationSettings/` is future-facing but should be planned now
- `emailQueue/` and `discordQueue/` can be introduced later without changing the core notification schema

---

## 21. Error Handling

### Notification creation failures

Log at minimum:

- notification type
- source event key
- recipient count
- team/league context
- route destination

Do not silently drop failures.

### Forum creation failures

Log at minimum:

- external event key
- forum id
- intended thread title
- generated thread/post IDs
- whether a duplicate check occurred

Retries must remain idempotent.

### Partial failure concerns

When writing many notifications:

- use Firestore batch writes where appropriate
- chunk large batches
- fail loudly on repository errors
- keep event keys available for replay

---

## 22. Observability

At minimum, track:

- `notifications_created_total`
- `notifications_failed_total`
- `forum_threads_created_total`
- `forum_posts_created_total`
- `idempotent_skips_total`
- `recipient_resolution_failures_total`

Log fields should include:

- `sourceEventKey`
- `notificationType`
- `routeName`
- `path`
- `threadId`
- `postId`
- `league`
- `teamId`

---

## 23. Testing Strategy

### Unit tests

Test:

- route builders
- notification title/message builders
- recipient resolvers
- idempotency checks
- metadata generation

### Repository tests

Test:

- create notifications
- create thread + first post atomically
- query by external event key
- query by source event key

### Integration tests

Test flows such as:

- injury event -> notifications created
- recruit signs -> notifications point to recruiting page
- postgame event -> thread and first post created
- duplicate event key -> duplicate creation prevented

---

## 24. Example Go Skeleton

```go
type notificationService struct {
    notifications NotificationRepository
    recipients    RecipientResolver
    routes        RouteBuilder
    clock         func() time.Time
}

func (s *notificationService) NotifyRecruitSigned(ctx context.Context, input RecruitSignedNotificationInput) error {
    exists, err := s.notifications.ExistsBySourceEventKey(ctx, input.SourceEventKey, NotificationRecruitSigned)
    if err != nil {
        return err
    }
    if exists {
        return nil
    }

    recipientUIDs, err := s.recipients.ResolveRecruitSignedRecipients(ctx, input.TeamID)
    if err != nil {
        return err
    }

    route := s.routes.BuildRecruitingRoute(input.League, input.TeamID)
    now := s.clock()

    notifications := make([]Notification, 0, len(recipientUIDs))
    for _, uid := range recipientUIDs {
        notifications = append(notifications, Notification{
            RecipientUID: uid,
            Type:         NotificationRecruitSigned,
            Title:        fmt.Sprintf("Recruit Signed: %s", input.RecruitName),
            Message:      fmt.Sprintf("%s has signed with %s.", input.RecruitName, input.TeamName),
            Severity:     "success",
            IsRead:       false,
            CreatedAt:    now,
            League:       input.League,
            TeamID:       input.TeamID,
            EntityType:   "recruit",
            EntityID:     strconv.Itoa(int(input.RecruitID)),
            Route:        route,
            SourceEventKey: input.SourceEventKey,
            Delivery: NotificationDelivery{
                InAppCreated: true,
            },
            Metadata: map[string]interface{}{
                "recruitId":   input.RecruitID,
                "recruitName": input.RecruitName,
            },
        })
    }

    return s.notifications.CreateNotifications(ctx, notifications)
}
```

---

## 25. Future Email and Discord Extension

This design is intentionally compatible with future secondary delivery channels.

### Email later

A future processor can:

- read notification events or notifications
- check `userNotificationSettings/{uid}`
- queue outbound email documents in `emailQueue/`

### Discord later

A future Discord worker can:

- read from `discordQueue/`
- map users or channels
- send bot messages

### Important recommendation

Treat email and Discord as delivery channels, not as the primary notification source.
The source of truth should remain the Go domain event pipeline.

---

## 26. Suggested Implementation Phases

### Phase 1: In-App Notifications

- notification schema
- route builder package
- recipient resolver package
- notification service for 5 core event types
- frontend consumption of route metadata

### Phase 2: Forum Automation

- thread + first post batch creation
- postgame thread creation
- recruiting sync admin threads
- forum reply notification support

### Phase 3: Preferences

- `userNotificationSettings/{uid}`
- per-type opt in/out
- in-app vs email vs Discord toggles

### Phase 4: Secondary Delivery

- email queue and processor
- Discord queue and worker
- digest support

---

## 27. Recommended Final Decisions

- Use typed notifications with stable constants
- Include explicit route metadata in every notification
- Use `SourceEventKey` for notification idempotency
- Use `ExternalEventKey` for automated forum thread/post idempotency
- Keep Firestore access behind repository interfaces
- Use a centralized route builder package
- Create automated threads and first posts in a single batch
- Design now for future email/Discord fan-out, but do not block on implementing them first

---

## 28. Conclusion

A Go + Firebase implementation is a strong fit for both in-app notifications and automated forum content.

The most important design choice is this:

**the backend must own the event type, recipients, and route payload.**

Once the notification contract includes:

- stable type
- title/message
- metadata
- route object
- idempotency key

the frontend becomes much simpler and more reliable.

For your use case, the best next steps are:

1. define notification type enums
2. implement centralized route builders
3. build a notification service with recipient resolution and idempotency
4. add forum thread/post automation with external event keys
5. later extend the same pipeline into email and Discord

This gives you one clean backend-first notification and forum system instead of several disconnected features.
