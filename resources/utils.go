package resources

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Validates the topic string that is used in ACL profiles
func validateACLTopicException(topic string) error {

	if strings.ContainsRune(topic, ' ') {
		return fmt.Errorf("topic string MUST not contain any whitespaces")
	}

	// Ensure that a topic exception string ends with @mqtt or @smf to define the syntax used
	if !strings.HasSuffix(topic, "@mqtt") && !strings.HasSuffix(topic, "@smf") {
		return fmt.Errorf("topic string MUST end with either '@smf' or '@mqtt' to specify the used syntax")
	}

	if strings.Contains(topic, "//") {
		return fmt.Errorf("topic string MUST not contain empty levels like //")
	}

	if strings.HasPrefix(topic, "/") {
		return fmt.Errorf("topic string MUST not start with /")
	}

	return nil
}

// Converts a string of comma separated days to a string with comma separated id's
func dayNamesToIds(schedule string) string {
	days := []struct {
		name string
		id   int
	}{
		{"Sunday", 0},
		{"Monday", 1},
		{"Tuesday", 2},
		{"Wednesday", 3},
		{"Thursday", 4},
		{"Friday", 5},
		{"Saturday", 6},
	}

	for _, day := range days {
		schedule = strings.ReplaceAll(schedule, day.name, strconv.Itoa(day.id))
	}
	return schedule
}

// Compares two slices of strings. Returns strings that are in `compare`, but not in `ref`, returns
// strings that are in `ref`, but not in `compare`
func SliceDelta(ref *[]string, compare *[]string) ([]string, []string) {
	var news, olds []string
	// build a map from reference
	refmap := make(map[string]struct{})
	for _, v := range *ref {
		refmap[v] = struct{}{}
	}
	// loop over compare - each found entry can be ignored. new entries must be added to news
	for _, v := range *compare {
		if _, ok := refmap[v]; !ok {
			news = append(news, v)
		} else {
			// remove the entry from the map
			delete(refmap, v)
		}
	}
	// all entries that are left over in the map are olds
	for key := range refmap {
		olds = append(olds, key)
	}

	return news, olds
}

// Compares to maps of strings. Returns a list with new entries which are in ref but not in compare, a list with entries in both maps, but with different values,
// and map of entries that are not in ref, but in compare.
func StringMapDelta(ref map[string]string, comp map[string]string) (map[string]string, map[string]string, map[string]string) {
	news := make(map[string]string)
	differents := make(map[string]string)
	olds := make(map[string]string)

	for k, v := range comp {
		olds[k] = v
	}

	// loop over all entries in reference and check against compare
	for k, v := range ref {
		if curval, ok := comp[k]; !ok {
			// entry not found in compare -> is new
			news[k] = v
		} else {
			if v != curval {
				// entry found, but different value
				differents[k] = v
			}
			// else: entry found and same value
		}
		delete(olds, k)
	}
	return news, differents, olds
}

func strHashFunc(obj interface{}) int {
	h := fnv.New64()
	h.Write([]byte(obj.(string)))
	return int(h.Sum64())
}

// Converts an anonymous map of strings (map[string]interface{}) to a real map of strings (map[map]string)
func iSMapTosSMap(iSMap *map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range *iSMap {
		result[k] = v.(string)
	}
	return result
}

// Converts an anonymous slice of strings ([]interface{}) to a real slice of strings ([]string)
func iSliceTosSlice(iSlice *[]interface{}) []string {
	target := make([]string, len(*iSlice))
	for i, v := range *iSlice {
		target[i] = v.(string)
	}
	return target
}

// Converts a slice of strings to an anonymous slice of interfaces ([]interface{})
func sSliceToiSlice(sSlice *[]string) []interface{} {
	target := make([]interface{}, len(*sSlice))
	for i, v := range *sSlice {
		target[i] = v
	}
	return target
}

func sSliceTosSet(sSlice *[]string) *schema.Set {
	return schema.NewSet(strHashFunc, sSliceToiSlice(sSlice))
}
