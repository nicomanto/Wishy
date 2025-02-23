<!DOCTYPE html>
<html lang="it">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="icon" type="image/png" href="https://wishyicons.s3.eu-west-1.amazonaws.com/gift-box.png">
    <title>{{.Username}}'s Wish List</title>
    <style>
        /* General Styles */
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f4f4f4;
            color: #333;
            text-align: center;
        }

        .container {
            max-width: 600px;
            margin: 40px auto;
            padding: 20px;
            background: white;
            border-radius: 12px;
            box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
        }

        h1 {
            font-size: 28px;
            color: #2c3e50;
        }

        /* Last Update & Button Container */
        .header-bar {
            display: flex;
            justify-content: space-between;
            align-items: center;
            width: 100%;
            padding: 10px 0;
            border-bottom: 2px solid #ddd;
        }

        .last-update {
            font-size: 14px;
            color: #777;
        }

        /* Download Button */
        .download-btn {
            background-color: #2980b9;
            color: white;
            border: none;
            padding: 8px 14px;
            border-radius: 6px;
            cursor: pointer;
            font-size: 14px;
            font-weight: bold;
            transition: background 0.2s;
        }

        .download-btn:hover {
            background-color: #1f6690;
        }

        .category {
            background: #fff;
            margin: 20px 0;
            padding: 20px;
            border-radius: 12px;
            box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
            text-align: left;
        }

        h2 {
            font-size: 20px;
            color: #2980b9;
            border-bottom: 3px solid #2980b9;
            padding-bottom: 10px;
            margin-bottom: 15px;
        }

        ul {
            list-style: none;
            padding: 0;
        }

        .wish-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 12px;
            background: #fafafa;
            margin: 8px 0;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
            transition: transform 0.2s;
        }

        .wish-item:hover {
            transform: translateY(-2px);
        }

        .wish-item a {
            color: #3498db;
            font-weight: 600;
            text-decoration: none;
        }

        .wish-item a:hover {
            text-decoration: underline;
        }

        .preference {
            font-size: 18px;
            color: #f1c40f;
        }

        @media (max-width: 600px) {
            .container {
                padding: 10px;
            }

            .header-bar {
                flex-direction: column;
                align-items: center;
                gap: 8px;
            }

            .download-btn {
                width: 100%;
            }
        }
    </style>
</head>

<body>
    <div class="container">
        <h1>🎁 {{.Username}}'s Wish List</h1>

        <!-- Last update and download button -->
        <div class="header-bar">
            <div class="last-update">Last updated: {{.LastUpdate}}</div>
            <button class="download-btn" onclick="downloadPDF()">⬇ Download PDF</button>
        </div>

        {{range .Wishes}}
        <div class="category">
            <h2>{{.Cat}}</h2>
            <ul>
                {{range .Wishes}}
                <li class="wish-item">
                    <a href="{{.Link}}" target="_blank">{{.Name}}</a>
                    <span class="preference">
                        {{if eq .Preference 3}} ⭐⭐⭐ {{end}}
                        {{if eq .Preference 2}} ⭐⭐ {{end}}
                        {{if eq .Preference 1}} ⭐ {{end}}
                    </span>
                </li>
                {{end}}
            </ul>
        </div>
        {{end}}
    </div>
</body>
<script>
    async function downloadPDF() {
    let baseUrl = window.location.origin;
    let pathParts = window.location.pathname.split('/');
    let stage = pathParts[1]; // "dev" or "prod"

    let params = new URLSearchParams(window.location.search);
    let uid = params.get("uid");

    if (!uid) {
        alert("User ID not found!");
        return;
    }

    // API URL for PDF
    let pdfUrl = `${baseUrl}/${stage}/wishes/pdf?uid=${uid}`;

    try {
        const response = await fetch(pdfUrl, {
            method: "GET",
            headers: {
                "Accept": "application/pdf" // Force request to expect a PDF
            }
        });

        if (!response.ok) {
            throw new Error("Failed to fetch PDF");
        }

        // Extract filename from Content-Disposition header
        let filename = "wishlist.pdf"; // Default filename
        const contentDisposition = response.headers.get("Content-Disposition");
        if (contentDisposition) {
            const match = contentDisposition.match(/filename="?([^";]+)"?/);
            if (match && match[1]) {
                filename = match[1];
            }
        }

        // Convert response to a Blob
        const pdfBlob = await response.blob();

        // Create a download link
        const downloadLink = document.createElement("a");
        downloadLink.href = URL.createObjectURL(pdfBlob);
        downloadLink.download = filename
        document.body.appendChild(downloadLink);
        downloadLink.click();
        document.body.removeChild(downloadLink);
    } catch (error) {
        console.error("Error downloading PDF:", error);
        alert("Error downloading PDF. Please try again.");
    }
}
</script>
</html>
