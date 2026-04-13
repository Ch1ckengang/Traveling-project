package tour

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"travel-backend/domain"
	"unicode"
)

// TourService - Tầng chứa business logic cho Tour
// Hiện tại chỉ đơn giản là lấy dữ liệu,
// nhưng sau này có thể thêm logic: lọc tour, sắp xếp, phân trang, v.v.

// GetAllTours - Lấy danh sách tất cả các tour
func GetAllTours() ([]domain.Tour, error) {
	return FindAllTours()
}

// GetDomesticTours - Lấy danh sách tour du lịch Việt Nam
func GetDomesticTours() ([]domain.Tour, error) {
	return FindDomesticTours()
}

// TourFilter - Bộ lọc nghiệp vụ cho danh sách tour.
type TourFilter struct {
	Category string
	City     string
	Duration string
	Price    string
	Sort     string
}

// GetToursByFilter - Lấy tour theo category và lọc nghiệp vụ từ ô tìm kiếm.
func GetToursByFilter(filter TourFilter) ([]domain.Tour, error) {
	tours, err := FindToursByCategory(filter.Category)
	if err != nil {
		return nil, err
	}

	city := strings.TrimSpace(strings.ToLower(filter.City))
	duration := strings.ToLower(filter.Duration)
	price := strings.ToLower(filter.Price)
	sortBy := strings.ToLower(filter.Sort)

	filtered := make([]domain.Tour, 0, len(tours))
	category := strings.ToLower(strings.TrimSpace(filter.Category))
	for _, tour := range tours {
		if category == "international" && !isInternationalTour(tour) {
			continue
		}

		if city != "" && !matchesCity(tour, city) {
			continue
		}

		if duration != "" && duration != "all" && !matchesDuration(tour.Duration, duration) {
			continue
		}

		if price != "" && price != "all" && !matchesPrice(tour.Price, price) {
			continue
		}

		filtered = append(filtered, tour)
	}

	applySort(filtered, sortBy)

	return filtered, nil
}

// CreateToursIfEmpty - Seed tour mẫu khi hệ thống chưa có tour nào
func CreateToursIfEmpty() error {
	return seedToursIfEmpty()
}

func matchesCity(tour domain.Tour, cityQuery string) bool {
	location := strings.ToLower(tour.Location)
	name := strings.ToLower(tour.Name)
	return strings.Contains(location, cityQuery) || strings.Contains(name, cityQuery)
}

func matchesDuration(durationText, durationFilter string) bool {
	days := extractDays(durationText)
	if days <= 0 {
		return strings.Contains(strings.ToLower(durationText), durationFilter)
	}

	switch durationFilter {
	case "short":
		return days <= 3
	case "medium":
		return days >= 4 && days <= 5
	case "long":
		return days >= 6
	default:
		return true
	}
}

func matchesPrice(priceText, priceFilter string) bool {
	amount := extractPrice(priceText)
	if amount <= 0 {
		return true
	}

	switch priceFilter {
	case "low":
		return amount < 3000000
	case "mid":
		return amount >= 3000000 && amount <= 7000000
	case "high":
		return amount > 7000000
	default:
		return true
	}
}

func extractDays(text string) int {
	re := regexp.MustCompile(`\d+`)
	value := re.FindString(text)
	if value == "" {
		return 0
	}

	days, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return days
}

func extractPrice(text string) int {
	b := strings.Builder{}
	for _, r := range text {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}

	if b.Len() == 0 {
		return 0
	}

	amount, err := strconv.Atoi(b.String())
	if err != nil {
		return 0
	}

	return amount
}

func applySort(tours []domain.Tour, sortBy string) {
	switch sortBy {
	case "price_asc":
		sort.SliceStable(tours, func(i, j int) bool {
			return extractPrice(tours[i].Price) < extractPrice(tours[j].Price)
		})
	case "price_desc":
		sort.SliceStable(tours, func(i, j int) bool {
			return extractPrice(tours[i].Price) > extractPrice(tours[j].Price)
		})
	case "duration_asc":
		sort.SliceStable(tours, func(i, j int) bool {
			return compareDuration(tours[i].Duration, tours[j].Duration) < 0
		})
	case "duration_desc":
		sort.SliceStable(tours, func(i, j int) bool {
			return compareDuration(tours[i].Duration, tours[j].Duration) > 0
		})
	case "name_asc":
		sort.SliceStable(tours, func(i, j int) bool {
			return strings.ToLower(tours[i].Name) < strings.ToLower(tours[j].Name)
		})
	case "name_desc":
		sort.SliceStable(tours, func(i, j int) bool {
			return strings.ToLower(tours[i].Name) > strings.ToLower(tours[j].Name)
		})
	case "latest":
		sort.SliceStable(tours, func(i, j int) bool {
			if tours[i].UpdatedAt.Equal(tours[j].UpdatedAt) {
				return tours[i].ID > tours[j].ID
			}
			return tours[i].UpdatedAt.After(tours[j].UpdatedAt)
		})
	}
}

func isInternationalTour(tour domain.Tour) bool {
	typeValue := strings.ToLower(strings.TrimSpace(tour.Type))
	countryNormalized := normalizeCountry(tour.Country)
	if typeValue != "international" {
		return false
	}

	if countryNormalized == "" {
		return false
	}

	return !isVietnamCountry(countryNormalized)
}

func isVietnamCountry(countryNormalized string) bool {
	switch countryNormalized {
	case "vietnam", "việtnam":
		return true
	default:
		return false
	}
}

func normalizeCountry(country string) string {
	normalized := strings.ToLower(strings.TrimSpace(country))
	normalized = strings.ReplaceAll(normalized, " ", "")
	return normalized
}

func compareDuration(left, right string) int {
	leftDays := extractDays(left)
	rightDays := extractDays(right)

	if leftDays == rightDays {
		return 0
	}

	if leftDays < rightDays {
		return -1
	}

	return 1
}
