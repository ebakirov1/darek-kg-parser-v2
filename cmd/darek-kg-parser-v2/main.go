package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/ebakirov1/darek-kg-parser-v2/internal/entity"
)

const (
	BaseURL = "http://address.darek.kg"
)

func main() {
	var (
		totalAreas   []entity.Area
		totalStreets []entity.Street
	)

	parentAreas := []int64{4948}

	for len(parentAreas) > 0 {
		newParentAreas := []int64{}

		for key, areaID := range parentAreas {
			var streets []entity.Street

			if err := queryStreets(&streets, areaID); err != nil {
				log.Fatal(err)
			}

			for _, street := range streets {
				street.AreaID = areaID
				totalStreets = append(totalStreets, street)
			}

			var areaResponse entity.AreaResponse

			if err := queryChild(&areaResponse, areaID); err != nil {
				log.Fatal(err)
			}

			for _, area := range areaResponse.Children {
				area.ParentID = areaID
				totalAreas = append(totalAreas, area)

				newParentAreas = append(newParentAreas, area.ID)
			}

			fmt.Println(key, len(parentAreas))
		}

		parentAreas = newParentAreas
	}

	if err := writeAreas(totalAreas); err != nil {
		log.Fatal(err)
	}

	if err := writeStreets(totalStreets); err != nil {
		log.Fatal(err)
	}
}

func queryStreets(streets *[]entity.Street, areaID int64) error {
	formData := url.Values{}
	formData.Set("ate_id", strconv.FormatInt(areaID, 10))

	if err := httpQuery(streets, "/ajax/street", formData); err != nil {
		return err
	}

	return nil
}

func queryChild(areaResponse *entity.AreaResponse, areaID int64) error {
	formData := url.Values{}
	formData.Set("ate", strconv.FormatInt(areaID, 10))

	if err := httpQuery(areaResponse, "/ajax/ateChild", formData); err != nil {
		return err
	}

	return nil
}

func writeStreets(totalStreets []entity.Street) error {
	f, err := os.OpenFile("data/streets.csv", os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer f.Close()

	csvWriter := csv.NewWriter(f)
	defer csvWriter.Flush()

	err = csvWriter.Write([]string{
		"area_id",
		"street_id",
		"street_type",
		"street_name",
		"street_name_full",
	})

	if err != nil {
		return err
	}

	for _, street := range totalStreets {
		err := csvWriter.Write([]string{
			strconv.FormatInt(street.AreaID, 10),
			strconv.FormatInt(street.StreetID, 10),
			strconv.FormatInt(street.Type, 10),
			street.Name,
			street.NameTP,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func writeAreas(totalAreas []entity.Area) error {
	f, err := os.OpenFile("data/areas.csv", os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer f.Close()

	csvWriter := csv.NewWriter(f)
	defer csvWriter.Flush()

	err = csvWriter.Write([]string{
		"parent_id",
		"area_id",
		"area_type",
		"area_type_name",
		"area_name",
	})

	if err != nil {
		return err
	}

	for _, ate := range totalAreas {
		err := csvWriter.Write([]string{
			strconv.FormatInt(ate.ParentID, 10),
			strconv.FormatInt(ate.ID, 10),
			strconv.FormatInt(ate.Type, 10),
			ate.TypeName,
			ate.Name,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func httpQuery(input any, urlPath string, data url.Values) error {
	resp, err := http.PostForm(BaseURL+urlPath, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(input); err != nil {
		return err
	}

	return nil
}
