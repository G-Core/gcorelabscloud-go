package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/lifecyclepolicy/v1/lifecyclepolicy"
	"github.com/G-Core/gcorelabscloud-go/gcore/schedule/v1/schedules"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func prepareResourceTestURL(scheduleID string) string {
	return fmt.Sprintf("/v1/schedule/%d/%d/%s", fake.ProjectID, fake.RegionID, scheduleID)
}

func TestGetCron(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareResourceTestURL(ScheduleCron1.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponseCron)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("schedule", "v1")
	ct, err := schedules.Get(client, ScheduleCron1.ID).Extract()
	require.NoError(t, err)

	sc, err := ct.Cook()
	require.NoError(t, err)
	cronSchedule, ok := sc.(lifecyclepolicy.CronSchedule)
	require.Equal(t, ok, true)

	require.NoError(t, err)
	require.Equal(t, ScheduleCron1, cronSchedule)
}

func TestGetInterval(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareResourceTestURL(ScheduleCron1.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponseInterval)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("schedule", "v1")
	ct, err := schedules.Get(client, ScheduleInterval1.ID).Extract()
	require.NoError(t, err)

	sc, err := ct.Cook()
	require.NoError(t, err)
	intervalSchedule, ok := sc.(lifecyclepolicy.IntervalSchedule)
	require.Equal(t, ok, true)

	require.NoError(t, err)
	require.Equal(t, ScheduleInterval1, intervalSchedule)
}

func TestUpdate(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareResourceTestURL(ScheduleCron1.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponseCron)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("schedule", "v1")

	opts := schedules.UpdateOpts{
		MaxQuantity:          2,
		DayOfWeek:            "fri, tue",
		Hour:                 "0, 20",
		Minute:               "30",
		Type:                 lifecyclepolicy.ScheduleTypeCron,
		ResourceNameTemplate: "reserve snap of the volume {volume_id}",
		Timezone:             "Europe/London",
		RetentionTime: lifecyclepolicy.RetentionTimer{
			Weeks: 2,
		},
	}

	ct, err := schedules.Update(client, ScheduleCron1.ID, opts).Extract()
	require.NoError(t, err)
	sc, err := ct.Cook()
	require.NoError(t, err)
	cronSchedule, ok := sc.(lifecyclepolicy.CronSchedule)
	require.Equal(t, ok, true)

	require.NoError(t, err)
	require.Equal(t, ScheduleCron1, cronSchedule)
}
