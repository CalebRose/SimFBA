package util

func GetRecruitingTendency(mod float64) string {
	/* translate this code from JS to GO
	   if (mod === 0.02) return "Signing Day Decision";
	   if (mod > 1.79) {
	     return "Will Sign Very Early";
	   } else if (mod < 1.8 && mod > 1.51) {
	     return "Will Sign Early";
	   } else if (mod < 1.5 && mod > 1.2) {
	     return "Average";
	   } else if (mod < 1.21 && mod > 1) {
	     return "Will Sign Late";
	   } else {
	     return "Will Sign Very Late";
	   }

	*/

	if mod == 0.02 {
		return "Signing Day Decision"
	}
	if mod > 1.79 {
		return "Will Sign Very Early"
	} else if mod < 1.8 && mod > 1.51 {
		return "Will Sign Early"
	} else if mod < 1.5 && mod > 1.2 {
		return "Average"
	} else if mod < 1.21 && mod > 1 {
		return "Will Sign Late"
	}
	return "Will Sign Very Late"
}
