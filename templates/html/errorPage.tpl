<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="icon" type="image/png" href="https://wishyicons.s3.eu-west-1.amazonaws.com/gift-box.png">
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;600;800&display=swap" rel="stylesheet">
    <title>{{.ErrorCode}} - {{.ErrorMessage}}</title>
    <style>
        :root {
            --bg: #0f172a;
            --card-bg: rgba(30, 41, 59, 0.7);
            --text-main: #f1f5f9;
            --text-muted: #94a3b8;
            --error-accent: #fb7185; /* Soft Rose/Red */
            --primary: #818cf8;
            --border: rgba(255, 255, 255, 0.08);
        }

        body {
            font-family: 'Inter', sans-serif;
            background-color: var(--bg);
            background-image: 
                radial-gradient(at 0% 0%, rgba(251, 113, 133, 0.1) 0px, transparent 50%),
                radial-gradient(at 100% 100%, rgba(99, 102, 241, 0.05) 0px, transparent 50%);
            color: var(--text-main);
            margin: 0;
            display: flex;
            align-items: center;
            justify-content: center;
            height: 100vh;
            padding: 20px;
            box-sizing: border-box;
        }

        .container {
            background: var(--card-bg);
            backdrop-filter: blur(12px);
            padding: 60px 40px;
            max-width: 500px;
            width: 100%;
            text-align: center;
            border-radius: 24px;
            border: 1px solid var(--border);
            box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
            animation: slideUp 0.6s cubic-bezier(0.16, 1, 0.3, 1);
        }

        .error-code {
            font-size: 100px;
            font-weight: 800;
            margin: 0;
            line-height: 1;
            background: linear-gradient(to bottom, #fff, var(--error-accent));
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            filter: drop-shadow(0 0 15px rgba(251, 113, 133, 0.3));
        }

        .error-message {
            font-size: 1.25rem;
            color: var(--text-main);
            margin-top: 20px;
            font-weight: 600;
        }

        .friendly-message {
            font-size: 1rem;
            color: var(--text-muted);
            margin-top: 10px;
            margin-bottom: 30px;
            line-height: 1.6;
        }

        /* Go Back Button */
        .btn {
            display: inline-flex;
            align-items: center;
            gap: 8px;
            padding: 12px 24px;
            font-size: 0.9rem;
            font-weight: 700;
            color: #0f172a;
            background: var(--primary);
            text-decoration: none;
            border-radius: 12px;
            transition: all 0.3s ease;
        }

        .btn:hover {
            background: #a5b4fc;
            transform: translateY(-2px);
            box-shadow: 0 10px 15px -3px rgba(129, 140, 248, 0.4);
        }

        @keyframes slideUp {
            from {
                opacity: 0;
                transform: translateY(20px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }

        @media (max-width: 480px) {
            .container { padding: 40px 20px; }
            .error-code { font-size: 80px; }
        }
    </style>
</head>
<body>

    <div class="container">
        <div class="error-code">{{.ErrorCode}}</div>
        <p class="error-message">Error Encountered</p>
        <p class="friendly-message">Oops! {{.FriendlyErrorMessage}}</p>
        
        <!-- Optional Go Back Button -->
        <a href="javascript:history.back()" class="btn">
            <svg width="18" height="18" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24"><path d="M10 19l-7-7m0 0l7-7m-7 7h18"></path></svg>
            Return to List
        </a>
    </div>

</body>
</html>