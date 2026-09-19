package utils

import (
"database/sql"
"net/http"
"encoding/json"
"log")

// The Scraper

type Category struct {
	Key       string  `json:"key"`
    ParentKey *string `json:"parentKey"`
}

type Offers struct {
    Data []Job `json:"data"`
}

type Job struct {
    Guid            string `json:"guid"`
    WorkplaceType   string `json:"workplaceType"`
    ExperienceLevel string `json:"experienceLevel"`
    City            string `json:"city"`
    ApplyUrl        string `json:"applyUrl"`
	Category 	  Category `json:"category"`
}

func Scrape(db *sql.DB) error {
	rows, err := db.Query("SELECT DISTINCT category FROM subscriptions")
	if err != nil {
		return err
	}
	defer rows.Close()
	var categories[]string
	for rows.Next() {
    var category string
    err = rows.Scan(&category)
    if err != nil {
        return err
    }
    categories = append(categories, category)
	}
	log.Println("Scraping jobs for categories:", categories)
	var incomeJobs []Job
	for i := 0; i < len(categories); i++ {
		response, err := http.Get("https://justjoin.it/api/candidate-api/offers?city=Warszawa&cityRadius=20&categories=" + categories[i] + "&sortBy=publishedAt&orderBy=descending&isPromoted=false")
		if err != nil {
			return err
		}
		defer response.Body.Close()
		log.Println(response.Status)
		var offer Offers
		err = json.NewDecoder(response.Body).Decode(&offer)
		if err != nil {
			return err
		}
		incomeJobs = append(incomeJobs, offer.Data...)
	}
	log.Println("Scraping completed. Inserting jobs into database...")
	log.Println("Number of jobs scraped:", len(incomeJobs))
	for i := 0; i < len(incomeJobs); i++ {
		_, err = db.Exec("INSERT INTO job_offers (guid, category, workplace, experience, location, applyurl) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (guid) DO NOTHING", incomeJobs[i].Guid, incomeJobs[i].Category.Key, incomeJobs[i].WorkplaceType, incomeJobs[i].ExperienceLevel, incomeJobs[i].City, incomeJobs[i].ApplyUrl)
		if err != nil {
			return err
		}
	}
	log.Println("Jobs inserted into database successfully.")
	return nil
}