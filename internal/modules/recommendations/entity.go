package recommendations

import (
	"github.com/karirnusantara/api/internal/modules/cvs"
	"github.com/karirnusantara/api/internal/modules/jobs"
)

// RecommendationScore represents a job with match score
type RecommendationScore struct {
	Job             *jobs.JobResponse `json:"job"`
	Score           int               `json:"score"`
	MatchReasons    []string          `json:"match_reasons"`
	MismatchReasons []string          `json:"mismatch_reasons,omitempty"`
}

// UserProfile represents the user's profile for matching
type UserProfile struct {
	UserID             uint64   `json:"user_id"`
	Name               string   `json:"name"`
	Skills             []string `json:"skills"`
	PreferredJobTypes  []string `json:"preferred_job_types"`
	PreferredLocations []string `json:"preferred_locations"`
	ExpectedSalaryMin  int64    `json:"expected_salary_min"`
	ExpectedSalaryMax  int64    `json:"expected_salary_max"`
	TotalExperience    float64  `json:"total_experience"` // in years
	ExperienceLevel    string   `json:"experience_level"` // entry, junior, mid, senior, lead
	Location           string   `json:"location"`         // current province
}

// RecommendationsResponse is the API response
type RecommendationsResponse struct {
	Recommendations []RecommendationScore `json:"recommendations"`
	TotalJobs       int                   `json:"total_jobs"`
	MatchedJobs     int                   `json:"matched_jobs"`
	AverageScore    int                   `json:"average_score"`
	ProfileComplete bool                  `json:"profile_complete"`
}

// BuildUserProfile builds UserProfile from CV and Profile data
func BuildUserProfile(userID uint64, name string, cv *cvs.CVResponse, province string, preferredLocations []string, preferredJobTypes []string, expectedSalaryMin, expectedSalaryMax int64) *UserProfile {
	profile := &UserProfile{
		UserID:             userID,
		Name:               name,
		Skills:             []string{},
		PreferredJobTypes:  preferredJobTypes,
		PreferredLocations: preferredLocations,
		ExpectedSalaryMin:  expectedSalaryMin,
		ExpectedSalaryMax:  expectedSalaryMax,
		TotalExperience:    0,
		Location:           province,
	}

	if cv != nil {
		// Extract skills
		for _, skill := range cv.Skills {
			profile.Skills = append(profile.Skills, skill.Name)
		}

		// Calculate total experience
		profile.TotalExperience = calculateTotalExperience(cv.Experience)
		profile.ExperienceLevel = determineExperienceLevel(profile.TotalExperience)
	}

	return profile
}

// calculateTotalExperience calculates total years of experience from CV
func calculateTotalExperience(experiences []cvs.Experience) float64 {
	if len(experiences) == 0 {
		return 0
	}

	totalMonths := 0
	for _, exp := range experiences {
		months := calculateMonthsBetween(exp.StartDate, exp.EndDate, exp.IsCurrent)
		totalMonths += months
	}

	return float64(totalMonths) / 12.0
}

// calculateMonthsBetween calculates months between two dates
func calculateMonthsBetween(startDate, endDate string, isCurrent bool) int {
	// Parse dates (format: YYYY-MM-DD or YYYY-MM)
	startYear, startMonth := parseYearMonth(startDate)

	var endYear, endMonth int
	if isCurrent || endDate == "" {
		// Use current date
		now := 2026 // Hardcoded for this context
		endYear = now
		endMonth = 1
	} else {
		endYear, endMonth = parseYearMonth(endDate)
	}

	months := (endYear-startYear)*12 + (endMonth - startMonth)
	if months < 0 {
		return 0
	}
	return months
}

// parseYearMonth extracts year and month from date string
func parseYearMonth(date string) (int, int) {
	if len(date) < 7 {
		return 0, 0
	}

	year := 0
	month := 0

	// Parse year
	for i := 0; i < 4 && i < len(date); i++ {
		year = year*10 + int(date[i]-'0')
	}

	// Parse month
	if len(date) >= 7 {
		month = int(date[5]-'0')*10 + int(date[6]-'0')
	}

	return year, month
}

// determineExperienceLevel determines level based on years
func determineExperienceLevel(years float64) string {
	switch {
	case years < 1:
		return "entry"
	case years < 2:
		return "junior"
	case years < 5:
		return "mid"
	case years < 8:
		return "senior"
	default:
		return "lead"
	}
}
