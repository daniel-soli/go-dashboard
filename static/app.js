document.getElementById("loadData").addEventListener("click", async function() {
    try {
        const response = await fetch("/api/data");
        const data = await response.json();
        document.getElementById("message").innerText = data.message;
    } catch (error) {
        document.getElementById("message").innerText = "Error fetching data.";
        console.error(error);
    }
});
