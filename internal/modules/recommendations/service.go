package recommendations

import (
	"strings"

	"github.com/karirnusantara/api/internal/modules/jobs"
)

// Service handles recommendation business logic
type Service struct{}

// NewService creates a new recommendations service
func NewService() *Service {
	return &Service{}
}

// CalculateJobScore calculates match score between user profile and job
// Weights:
// - Skills Match: 35%
// - Experience Level: 25%
// - Location Match: 20%
// - Salary Match: 10%
// - Job Type Match: 10%
func (s *Service) CalculateJobScore(profile *UserProfile, job *jobs.JobResponse) RecommendationScore {
	var score float64 = 0
	matchReasons := []string{}
	mismatchReasons := []string{}

	// 1. Skills Match (35%)
	skillScore, skillMatches := s.calculateSkillScore(profile.Skills, job)
	score += skillScore * 0.35
	if len(skillMatches) > 0 {
		matchReasons = append(matchReasons, "Skill cocok: "+strings.Join(skillMatches[:min(3, len(skillMatches))], ", "))
	} else if len(profile.Skills) > 0 {
		mismatchReasons = append(mismatchReasons, "Skill Anda tidak sesuai dengan requirements")
	}

	// 2. Experience Level Match (25%)
	expScore := s.calculateExperienceScore(profile.ExperienceLevel, profile.TotalExperience, job.ExperienceLevel)
	score += expScore * 0.25
	if expScore >= 80 {
		matchReasons = append(matchReasons, "Pengalaman Anda sesuai")
	} else if expScore >= 50 {
		matchReasons = append(matchReasons, "Pengalaman Anda mendekati")
	} else if profile.TotalExperience > 0 {
		mismatchReasons = append(mismatchReasons, "Level pengalaman berbeda")
	}

	// 3. Location Match (20%)
	locScore := s.calculateLocationScore(profile, job)
	score += locScore * 0.20
	if locScore >= 100 {
		if job.Location.IsRemote {
			matchReasons = append(matchReasons, "Remote - Bisa kerja dari mana saja")
		} else {
			matchReasons = append(matchReasons, "Lokasi sesuai: "+job.Location.Province)
		}
	} else if locScore >= 50 {
		matchReasons = append(matchReasons, "Lokasi mendekati preferensi")
	}

	// 4. Salary Match (10%)
	salScore := s.calculateSalaryScore(profile, job)
	score += salScore * 0.10
	if salScore >= 80 && job.Salary != nil && (job.Salary.Min > 0 || job.Salary.Max > 0) {
		matchReasons = append(matchReasons, "Gaji sesuai ekspektasi")
	}

	// 5. Job Type Match (10%)
	typeScore := s.calculateJobTypeScore(profile.PreferredJobTypes, job.JobType)
	score += typeScore * 0.10
	if typeScore >= 100 {
		matchReasons = append(matchReasons, "Tipe pekerjaan sesuai")
	}

	// Bonus: Remote option
	if job.Location.IsRemote && !containsRemote(matchReasons) {
		score = minFloat(score+5, 100)
	}

	return RecommendationScore{
		Job:             job,
		Score:           int(score),
		MatchReasons:    matchReasons[:min(4, len(matchReasons))],
		MismatchReasons: mismatchReasons[:min(2, len(mismatchReasons))],
	}
}

// calculateSkillScore scores skill matches (0-100)
func (s *Service) calculateSkillScore(userSkills []string, job *jobs.JobResponse) (float64, []string) {
	if len(userSkills) == 0 {
		return 50, nil // Default score if no skills provided
	}

	matched := []string{}
	searchText := strings.ToLower(job.Title + " " + job.Description)
	if job.Requirements != "" {
		searchText += " " + strings.ToLower(job.Requirements)
	}

	// Common tech skill aliases
	aliases := map[string][]string{
		"javascript": {"js", "javascript", "node.js", "nodejs"},
		"typescript": {"ts", "typescript"},
		"react":      {"react", "reactjs", "react.js"},
		"python":     {"python", "django", "flask"},
		"go":         {"go", "golang"},
		"sql":        {"sql", "mysql", "postgresql", "postgres"},
		"java":       {"java", "spring"},
	}

	for _, skill := range userSkills {
		skillLower := strings.ToLower(skill)

		// Check direct match
		if strings.Contains(searchText, skillLower) {
			matched = append(matched, skill)
			continue
		}

		// Check aliases
		for base, aliasList := range aliases {
			if skillLower == base || contains(aliasList, skillLower) {
				for _, alias := range aliasList {
					if strings.Contains(searchText, alias) {
						matched = append(matched, skill)
						break
					}
				}
			}
		}
	}

	if len(matched) == 0 {
		return 20, nil
	}

	// Score based on match ratio
	matchRatio := float64(len(matched)) / float64(len(userSkills))
	score := 50 + (matchRatio * 50) // 50-100 based on match ratio

	return score, matched
}

// calculateExperienceScore scores experience level match (0-100)
func (s *Service) calculateExperienceScore(userLevel string, userYears float64, jobLevel string) float64 {
	if userLevel == "" || userYears == 0 {
		return 50 // Default if no experience data
	}

	levels := map[string]int{
		"entry":     0,
		"junior":    1,
		"mid":       2,
		"senior":    3,
		"lead":      4,
		"executive": 5,
	}

	userLevelNum := levels[userLevel]
	jobLevelNum := levels[jobLevel]

	diff := userLevelNum - jobLevelNum

	switch {
	case diff == 0:
		return 100
	case diff == 1:
		return 85 // Overqualified by 1 level
	case diff == -1:
		return 70 // Underqualified by 1 level
	case diff > 1:
		return 60 // Very overqualified
	default:
		return 40 // Very underqualified
	}
}

// calculateLocationScore scores location match (0-100)
func (s *Service) calculateLocationScore(profile *UserProfile, job *jobs.JobResponse) float64 {
	// Remote jobs always match
	if job.Location.IsRemote {
		return 100
	}

	jobProvince := strings.ToLower(job.Location.Province)
	jobCity := strings.ToLower(job.Location.City)

	// Check preferred locations
	for _, loc := range profile.PreferredLocations {
		locLower := strings.ToLower(loc)
		if strings.Contains(jobProvince, locLower) || strings.Contains(locLower, jobProvince) {
			return 100
		}
		if strings.Contains(jobCity, locLower) || strings.Contains(locLower, jobCity) {
			return 90
		}
	}

	// Check current location
	if profile.Location != "" {
		profileLoc := strings.ToLower(profile.Location)
		if strings.Contains(jobProvince, profileLoc) || strings.Contains(profileLoc, jobProvince) {
			return 80
		}
	}

	return 30 // No location match
}

// calculateSalaryScore scores salary match (0-100)
func (s *Service) calculateSalaryScore(profile *UserProfile, job *jobs.JobResponse) float64 {
	if job.Salary == nil || (job.Salary.Min == 0 && job.Salary.Max == 0) {
		return 50 // No salary info, neutral score
	}

	if profile.ExpectedSalaryMin == 0 && profile.ExpectedSalaryMax == 0 {
		return 70 // User has no preference, slightly positive
	}

	jobMin := job.Salary.Min
	jobMax := job.Salary.Max

	// Check if ranges overlap
	if profile.ExpectedSalaryMax > 0 && jobMin > 0 {
		if profile.ExpectedSalaryMax < jobMin {
			return 30 // Job pays more than expected max
		}
	}

	if profile.ExpectedSalaryMin > 0 && jobMax > 0 {
		if profile.ExpectedSalaryMin > jobMax {
			return 20 // Job pays less than expected min
		}
	}

	// Ranges overlap
	return 100
}

// calculateJobTypeScore scores job type match (0-100)
func (s *Service) calculateJobTypeScore(preferredTypes []string, jobType string) float64 {
	if len(preferredTypes) == 0 {
		return 70 // No preference, neutral-positive
	}

	// Map API job types to human-readable
	typeMap := map[string][]string{
		"full_time":  {"full_time", "full-time", "fulltime"},
		"part_time":  {"part_time", "part-time", "parttime"},
		"contract":   {"contract", "kontrak"},
		"internship": {"internship", "magang", "intern"},
		"freelance":  {"freelance"},
	}

	aliases := typeMap[jobType]
	for _, pref := range preferredTypes {
		prefLower := strings.ToLower(pref)
		if prefLower == jobType || contains(aliases, prefLower) {
			return 100
		}
	}

	return 30
}

// GetRecommendations returns recommended jobs for user
func (s *Service) GetRecommendations(profile *UserProfile, allJobs []*jobs.JobResponse, limit int) *RecommendationsResponse {
	if limit <= 0 {
		limit = 20
	}

	recommendations := []RecommendationScore{}

	for _, job := range allJobs {
		rec := s.CalculateJobScore(profile, job)
		if rec.Score >= 30 { // Only include jobs with >= 30% match
			recommendations = append(recommendations, rec)
		}
	}

	// Sort by score descending
	for i := 0; i < len(recommendations)-1; i++ {
		for j := i + 1; j < len(recommendations); j++ {
			if recommendations[j].Score > recommendations[i].Score {
				recommendations[i], recommendations[j] = recommendations[j], recommendations[i]
			}
		}
	}

	// Limit results
	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	// Calculate average score
	avgScore := 0
	if len(recommendations) > 0 {
		total := 0
		for _, r := range recommendations {
			total += r.Score
		}
		avgScore = total / len(recommendations)
	}

	// Profile completeness check
	profileComplete := len(profile.Skills) > 0 || profile.TotalExperience > 0 || len(profile.PreferredLocations) > 0

	return &RecommendationsResponse{
		Recommendations: recommendations,
		TotalJobs:       len(allJobs),
		MatchedJobs:     len(recommendations),
		AverageScore:    avgScore,
		ProfileComplete: profileComplete,
	}
}

// Helper functions
func contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func containsRemote(reasons []string) bool {
	for _, r := range reasons {
		if strings.Contains(strings.ToLower(r), "remote") {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
