package strava

import "fmt"

func GetAllActivities(accessToken string)([]Activity, error){

	var allActivities []Activity
	page := 1
	perPage := 100

	for{

		activities, err := GetActivitiesPage(accessToken, page, perPage)
		if err != nil{
			fmt.Println("Error fetching activities using GetActivitiesPage")
			return nil, err
		}
		
		//check len(activities)- if its 0 then stop 
		if len(activities) == 0{
			break
		}

		//activities ... beacuse it takes all elements from activities slice and append them individually
		//append() expects type Activity and not []Activity 
		allActivities = append(
			allActivities,
		    activities... ) 
		
		fmt.Printf(
			"Fetched page %d (%d activities)\n",
			page,
			len(activities),
		)

		page ++

	}

	return allActivities, nil
	
}