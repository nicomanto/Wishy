<!DOCTYPE html>
<html lang="it">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="icon" type="image/png" href="https://wishyicons.s3.eu-west-1.amazonaws.com/gift-box.png">
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;600;700&display=swap" rel="stylesheet">
    <title>{{.Username}}'s Wish List</title>
    <style>
        :root {
            /* Dark Theme Colors */
            --bg: #0f172a; /* Deep Slate */
            --card-bg: #1e293b; /* Slate 800 */
            --primary: #818cf8; /* Indigo 400 */
            --primary-hover: #a5b4fc;
            --accent: #fbbf24; /* Amber 400 */
            --text-main: #f1f5f9;
            --text-muted: #94a3b8;
            --border: rgba(255, 255, 255, 0.08);
        }

        body {
            font-family: 'Inter', sans-serif;
            margin: 0;
            padding: 0;
            background-color: var(--bg);
            /* Soft mesh gradient for depth */
            background-image: 
                radial-gradient(at 0% 0%, rgba(99, 102, 241, 0.15) 0px, transparent 50%),
                radial-gradient(at 100% 100%, rgba(30, 41, 59, 0.4) 0px, transparent 50%);
            color: var(--text-main);
            min-height: 100vh;
        }

        .container {
            max-width: 700px;
            margin: 0 auto;
            padding: 40px 20px;
        }

        header {
            text-align: center;
            margin-bottom: 40px;
        }

        h1 {
            font-size: 2.5rem;
            font-weight: 800;
            letter-spacing: -0.025em;
            margin-bottom: 15px;
            background: linear-gradient(to right, #fff, #94a3b8);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
        }

        /* Glassmorphism Header Bar */
        .header-bar {
            display: flex;
            justify-content: space-between;
            align-items: center;
            background: rgba(30, 41, 59, 0.5);
            backdrop-filter: blur(12px);
            padding: 14px 24px;
            border-radius: 20px;
            border: 1px solid var(--border);
            box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.3);
        }

        .last-update {
            font-size: 0.85rem;
            color: var(--text-muted);
            font-weight: 500;
        }

        /* Neon Style Download Button */
        .download-btn {
            background-color: var(--primary);
            color: #0f172a;
            border: none;
            padding: 10px 20px;
            border-radius: 12px;
            cursor: pointer;
            font-size: 0.875rem;
            font-weight: 700;
            transition: all 0.3s ease;
            display: flex;
            align-items: center;
            gap: 8px;
        }

        .download-btn:hover {
            background-color: var(--primary-hover);
            transform: translateY(-2px);
            box-shadow: 0 0 20px rgba(129, 140, 248, 0.4);
        }

        /* Categories */
        .category {
            margin-top: 40px;
        }

        h2 {
            font-size: 1.1rem;
            text-transform: uppercase;
            letter-spacing: 0.1em;
            color: var(--text-muted);
            margin-bottom: 20px;
            display: flex;
            align-items: center;
            gap: 15px;
        }

        h2::after {
            content: "";
            flex: 1;
            height: 1px;
            background: var(--border);
        }

        ul {
            list-style: none;
            padding: 0;
            margin: 0;
        }

        /* Dark Cards */
        .wish-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 18px 24px;
            background: var(--card-bg);
            margin-bottom: 12px;
            border-radius: 16px;
            border: 1px solid var(--border);
            transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        }

        .wish-item:hover {
            background: #243147;
            border-color: rgba(129, 140, 248, 0.3);
            transform: scale(1.02);
            box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.2);
        }

        .wish-info {
            display: flex;
            flex-direction: column;
        }

        .wish-item a {
            color: var(--text-main);
            font-weight: 600;
            text-decoration: none;
            font-size: 1.1rem;
            display: flex;
            align-items: center;
            gap: 6px;
        }

        .wish-item a:hover {
            color: var(--primary);
        }

        .preference {
            letter-spacing: 3px;
            font-size: 0.8rem;
            opacity: 0.9;
        }

        /* Pulsing Amber Signal for New items */
        .signal {
            border-left: 4px solid var(--accent) !important;
            background: linear-gradient(to right, rgba(251, 191, 36, 0.05), var(--card-bg));
            animation: dark-pulse 2s infinite;
        }

        .new-badge {
            font-size: 0.6rem;
            background: var(--accent);
            color: #000;
            padding: 2px 8px;
            border-radius: 20px;
            text-transform: uppercase;
            font-weight: 900;
            margin-bottom: 6px;
            width: fit-content;
        }

        @keyframes dark-pulse {
            0% { border-color: rgba(251, 191, 36, 0.4); }
            50% { border-color: rgba(251, 191, 36, 1); }
            100% { border-color: rgba(251, 191, 36, 0.4); }
        }

        /* Mobile Adjustments */
        @media (max-width: 600px) {
            .container { padding: 20px 15px; }
            h1 { font-size: 1.8rem; }
            .header-bar { flex-direction: column; gap: 15px; text-align: center; }
            .wish-item { flex-direction: column; align-items: flex-start; gap: 12px; }
            .preference { align-self: flex-end; }
        }
    </style>
</head>

<body>
    <div class="container">
        <header>
            <h1>🎁 {{.Username}}'s Wish List</h1>
            <div class="header-bar">
                <div class="last-update">Last Sync: {{.LastUpdate}}</div>
                <button class="download-btn" onclick="downloadPDF()">
                    <svg width="18" height="18" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24"><path d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"></path></svg>
                    PDF Export
                </button>
            </div>
        </header>

        {{range .Wishes}}
        <div class="category">
            <h2>{{.Cat}}</h2>
            <ul>
                {{range .Wishes}}
                <li class="wish-item {{if .IsRecent}}signal{{end}}">
                    <div class="wish-info">
                        {{if .IsRecent}}<span class="new-badge">New Arrival</span>{{end}}
                        <a href="{{.Link}}" target="_blank">
                            {{.Name}}
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="7" y1="17" x2="17" y2="7"></line><polyline points="7 7 17 7 17 17"></polyline></svg>
                        </a>
                    </div>
                    <span class="preference">
                        {{if eq .Preference 3}}⭐⭐⭐{{end}}
                        {{if eq .Preference 2}}⭐⭐{{end}}
                        {{if eq .Preference 1}}⭐{{end}}
                    </span>
                </li>
                {{end}}
            </ul>
        </div>
        {{end}}
    </div>

    <script>
        async function downloadPDF() {
            let baseUrl = window.location.origin;
            let pathParts = window.location.pathname.split('/');
            let stage = pathParts[1]; 
            let params = new URLSearchParams(window.location.search);
            let uid = params.get("uid");

            if (!uid) {
                alert("User ID missing");
                return;
            }

            let pdfUrl = `${baseUrl}/${stage}/wishes/pdf?uid=${uid}`;

            try {
                const response = await fetch(pdfUrl, {
                    method: "GET",
                    headers: { "Accept": "application/pdf" }
                });

                if (!response.ok) throw new Error("Network error");

                const pdfBlob = await response.blob();
                const downloadLink = document.createElement("a");
                downloadLink.href = URL.createObjectURL(pdfBlob);
                downloadLink.download = "wishlist.pdf";
                document.body.appendChild(downloadLink);
                downloadLink.click();
                document.body.removeChild(downloadLink);
            } catch (error) {
                console.error(error);
                alert("Failed to generate PDF.");
            }
        }
    </script>
</body>
</html>