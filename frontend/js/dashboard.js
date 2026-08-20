// ==========================================
// ACTIVITY EMOJIS
// ==========================================

const activityEmojis = {

    Run: "🏃",
    Swim: "🏊",
    WeightTraining: "🏋️",
    Ride: "🚴",
    Hike: "🥾",
    Walk: "🚶",
    Yoga: "🧘",

    default: "🏅"
};


// ==========================================
// METRIC EMOJIS
// ==========================================

const metricEmojis = {

    duration: "⏱️",
    distance: "〽️",
    pace: "🏃",
    speed: "📈",
    elevation: "⛰️",
    heartRate: "❤️",
    power: "⚡"
};

// ==========================================
// SPORT NAMES
// ==========================================

const sportNames = {

    Run: "Running",

    Hike: "Hiking",

    Walk: "Walking",

    WeightTraining: "Weight Training",

    Swim: "Swimming",

    Ride: "Cycling"

};

// ==========================================
// AI SUMMARY
// ==========================================

const summaryButton = document.getElementById("summary-button");
const summaryBox = document.getElementById("summary-box");

summaryButton.addEventListener("click", async function () {

    summaryBox.textContent = "💭 Generating your 30-day summary...";

    try {

        const response = await fetch("/coach/report");

        if (!response.ok) {
            throw new Error("Failed to fetch summary");
        }

        const data = await response.json();

        summaryBox.textContent = data.summary;

    } catch (error) {

        console.error("Error generating summary:", error);

        summaryBox.textContent =
            "Unable to generate the summary. Please try again.";
    }
});


// ==========================================
// TOP SPORT 
// ==========================================

const topSportButton = document.getElementById("top-sport-button");
const topSportBox = document.getElementById("top-sport-section");


topSportButton.addEventListener("click", async function(){

    topSportBox.textContent = "💭 Checking your top sport in the last 30 days...";

    try{

        const response = await fetch("/stats/top-sport");
        if (! response.ok){
            throw new Error("Failed to fetch top sport");

        }

        const data = await response.json();

        //fetch the list of top sport/ top sports

        const sports = data.sports;

        if (sports.length === 0) {

            topSportBox.textContent =
                "No activities found in the last 30 days.";

            return;
        }

        const formattedSports = sports.map(function (sport) {

            const sportEmoji =
                activityEmojis[sport.sport] || "🏅";

            const sportName =
                sportNames[sport.sport] || sport.sport;

            return `${sportName} ${sportEmoji}`;
        });

        const sportText = formattedSports.join(" and ");

        const activityCount = sports[0].count;

        topSportBox.innerHTML =
            `Your top sport${sports.length > 1 ? "s" : ""} ` +
            `in the last 30 days ${sports.length > 1 ? "are..." : "is..."} ` +
            `<br><strong>${sportText}</strong></br>` +
            `You logged <strong>${activityCount} activities`+
            `${sports.length >1 ? "each." : "."}` ;

    } catch(error){

        console.error("Error fetching top sport", error)
        topSportBox.textContent = "Error fetching top sport. Please try again"

    }
});

// ==========================================
// PERSONALIZED AI COACHING
// ==========================================

const coachingGoal = document.getElementById("coaching-goal");
const coachingButton = document.getElementById("coaching-button");
const coachingBox = document.getElementById("coaching-box");

coachingButton.addEventListener("click", async function () {
    const goal = coachingGoal.value.trim();

    if (!goal) {
        coachingBox.textContent =
            "Please describe your running goal first.";
        return;
    }

    coachingButton.disabled = true;
    coachingBox.textContent =
        "🤔Thinking about your goal and recent running data...";

    try {
        const response = await fetch("/coach/coaching", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                goal: goal
            })
        });

        if (!response.ok) {
            throw new Error("Failed to generate coaching");
        }

        const data = await response.json();

        renderCoachingResponse(data.coaching);

    } catch (error) {
        console.error("Error generating personalized coaching:", error);

        coachingBox.textContent =
            "Unable to generate personalized coaching. Please try again.";

    } finally {
        coachingButton.disabled = false;
    }
});

function renderCoachingResponse(coaching) {
    coachingBox.innerHTML = "";

    const sections = [
        {
            title: "Observations",
            values: coaching.observations
        },
        {
            title: "Progress",
            values: coaching.progress
        },
        {
            title: "Risks",
            values: coaching.risks
        },
        {
            title: "Recommendations",
            values: coaching.recommendations
        }
    ];

    sections.forEach(function (section) {
        if (
            !section.values ||
            section.values.length === 0
        ) {
            return;
        }

        const heading = document.createElement("h3");
        heading.textContent = section.title;

        const list = document.createElement("ul");

        section.values.forEach(function (value) {
            const item = document.createElement("li");
            item.textContent = value;
            list.appendChild(item);
        });

        coachingBox.appendChild(heading);
        coachingBox.appendChild(list);
    });
}


// ==========================================
// RECENT ACTIVITIES
// ==========================================

const activitiesList = document.getElementById("activities-list");
const activitiesFooter = document.getElementById("activities-footer");


// Which metrics should each activity type display?

const activityConfig = {

    Run: [
        "duration",
        "distance",
        "pace",
        "elevation",
        "heartRate",
        "power"
    ],

    Hike: [
        "duration",
        "distance",
        "speed",
        "elevation",
        "heartRate"
    ],

    Walk: [
        "duration",
        "distance",
        "speed",
        "elevation",
        "heartRate"
    ],

    Swim: [
        "duration",
        "distance",
        "speed",
        "heartRate"
    ],

    WeightTraining: [
        "duration",
        "heartRate"
    ],

    Ride: [
        "duration",
        "distance",
        "speed",
        "elevation",
        "heartRate",
        "power"
    ],

    default: [
        "duration",
        "distance",
        "speed",
        "heartRate"
    ]
};



// ==========================================
// FORMATTING FUNCTIONS
// ==========================================

function formatDuration(seconds) {

    if (!seconds || seconds <= 0) {
        return "N/A";
    }

    const hours = Math.floor(seconds / 3600);

    const minutes = Math.floor(
        (seconds % 3600) / 60
    );

    if (hours > 0) {
        return `${hours}h ${minutes}m`;
    }

    return `${minutes}m`;
}


function formatDistance(meters) {

    if (!meters || meters <= 0) {
        return "N/A";
    }

    return `${(meters / 1000).toFixed(1)} km`;
}


// Convert meters/second into minutes/km.
//
// Example:
//
// 1.9 m/s
// ↓
// 1900 m in 1000 seconds
// ↓
// approximately 8:46 min/km

//format pace to min/km only for runs
function formatPace(speedMetersPerSecond) {

    if (!speedMetersPerSecond || speedMetersPerSecond <= 0) {
        return "N/A";
    }

    const secondsPerKm =
        1000 / speedMetersPerSecond;

    const minutes =
        Math.floor(secondsPerKm / 60);

    const seconds =
        Math.round(secondsPerKm % 60);

    if (seconds === 60) {
        return `${minutes + 1}:00 min/km`;
    }

    return `${minutes}:${String(seconds).padStart(2, "0")} min/km`;
}

//format pace to km/hr for all activities except runs
function formatSpeed(speedMetersPerSecond) {

    if (!speedMetersPerSecond || speedMetersPerSecond <= 0) {
        return "N/A";
    }

    const speedKmPerHour =
        speedMetersPerSecond * 3.6;

    return `${speedKmPerHour.toFixed(1)} km/h`;
}

function formatDate(dateString) {

    if (!dateString) {
        return "Unknown date";
    }

    const date = new Date(dateString);

    return date.toLocaleDateString(
        "en-GB",
        {
            day: "numeric",
            month: "short",
            year: "numeric"
        }
    );
}


// ==========================================
// ACTIVITY TYPE
// ==========================================

function getActivityType(activity) {

    return activity.sport_type ||
           activity.type ||
           "Activity";
}


function getActivityEmoji(activity) {

    const activityType =
        getActivityType(activity);

    return activityEmojis[activityType] ||
           activityEmojis.default;
}


// ==========================================
// CREATE ONE METRIC
// ==========================================

function createMetric(label, value, emoji) {

    const metric =
        document.createElement("div");

    metric.className =
        "activity-metric";


    const metricHeader =
        document.createElement("div");

    metricHeader.className =
        "metric-header";


    const emojiElement =
        document.createElement("span");

    emojiElement.className =
        "metric-emoji";

    emojiElement.textContent =
        emoji;


    const labelElement =
        document.createElement("span");

    labelElement.className =
        "metric-label";

    labelElement.textContent =
        label;


    const valueElement =
        document.createElement("strong");

    valueElement.className =
        "metric-value";

    valueElement.textContent =
        value;


    metricHeader.appendChild(emojiElement);
    metricHeader.appendChild(labelElement);

    metric.appendChild(metricHeader);
    metric.appendChild(valueElement);

    return metric;
}


// ==========================================
// CREATE METRICS FOR ONE ACTIVITY
// ==========================================

function createActivityMetrics(activity) {

    const metricsContainer =
        document.createElement("div");

    metricsContainer.className =
        "activity-metrics";


    const activityType =
        getActivityType(activity);


    const metrics =
        activityConfig[activityType] ||
        activityConfig.default;


    // Duration

    if (metrics.includes("duration")) {

        metricsContainer.appendChild(
            createMetric(
                "Duration",
                formatDuration(activity.moving_time),
                metricEmojis.duration
            )
        );
    }


    // Distance

    if (metrics.includes("distance")) {

        metricsContainer.appendChild(
            createMetric(
                "Distance",
                formatDistance(activity.distance),
                metricEmojis.distance
            )
        );
    }


    // Pace

    if (metrics.includes("pace")) {

        metricsContainer.appendChild(
            createMetric(
                "Pace",
                formatPace(activity.average_speed),
                metricEmojis.pace
            )
        );
    }


    // Elevation

    if (metrics.includes("elevation")) {

        metricsContainer.appendChild(
            createMetric(
                "Elevation",
                `${Math.round(activity.total_elevation_gain || 0)} m`,
                metricEmojis.elevation
            )
        );
    }


    // Heart rate

    if (metrics.includes("heartRate")) {

        const heartRate =
            activity.average_heartrate > 0
                ? `${Math.round(activity.average_heartrate)} bpm`
                : "N/A";

        metricsContainer.appendChild(
            createMetric(
                "Avg HR",
                heartRate,
                metricEmojis.heartRate
            )
        );
    }


    // Power

    if (metrics.includes("power")) {

        const power =
            activity.average_watts > 0
                ? `${Math.round(activity.average_watts)} W`
                : "N/A";

        metricsContainer.appendChild(
            createMetric(
                "Avg Power",
                power,
                metricEmojis.power
            )
        );
    }

    // Speed

if (metrics.includes("speed")) {

    metricsContainer.appendChild(

        createMetric(

            "Avg Speed",

            formatSpeed(activity.average_speed),

            metricEmojis.speed

        )

    );

}


    return metricsContainer;
}


// ==========================================
// CREATE ONE ACTIVITY CARD
// ==========================================

function createActivityCard(activity) {

    const card =
        document.createElement("article");

    card.className =
        "activity-card";


    // Activity title

    const title =
        document.createElement("h3");

    title.textContent =
        activity.name || "Unnamed Activity";


    // Activity date

    const date =
        document.createElement("p");

    date.className =
        "activity-date";

    date.textContent =
        `📅 ${formatDate(activity.start_date_local)}`;


    // Activity type

    const type =
        document.createElement("div");

    type.className =
        "activity-type";

    type.textContent =
        `${getActivityEmoji(activity)} ${getActivityType(activity)}`;


    // Activity metrics

    const metrics =
        createActivityMetrics(activity);


    card.appendChild(title);
    card.appendChild(date);
    card.appendChild(type);
    card.appendChild(metrics);

    return card;
}


// ==========================================
// LOAD RECENT ACTIVITIES
// ==========================================

async function loadRecentActivities() {

    try {

        const response =
            await fetch("/activities/recent");


        if (!response.ok) {

            throw new Error(
                "Failed to fetch activities"
            );
        }


        const activities =
            await response.json();


        // Remove the loading message

        activitiesList.innerHTML = "";


        // Create one card for every activity

        activities.forEach(function (activity) {

            const card =
                createActivityCard(activity);

            activitiesList.appendChild(card);
        });


        // Footer

        activitiesFooter.textContent =
            `Showing your last ${activities.length} activities from Strava`;


    } catch (error) {

        console.error(
            "Error loading activities:",
            error
        );

        activitiesList.innerHTML =
            "<p>Unable to load recent activities.</p>";

        activitiesFooter.textContent = "";
    }
}


// ==========================================
// LOAD ACTIVITIES WHEN DASHBOARD OPENS
// ==========================================

loadRecentActivities();