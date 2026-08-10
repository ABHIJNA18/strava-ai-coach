const summaryButton = document.getElementById('summary-button');
const summaryBox = document.getElementById('summary-box');

summaryButton.addEventListener("click", async function(){

    summaryBox.textContent = "Generating your 30-day summary..."

    try {
        const response = await fetch("/coach/report");

        if (!response.ok){

            throw new Error("Failed to fetch summary")
        }

        const data = await response.json();
        summaryBox.textContent = data.summary;

    } catch(error) {
        console.error("Error generating summary:", error);
        summaryBox.textContent = "Unable to generate the summary. Please try again.";

    }

});
