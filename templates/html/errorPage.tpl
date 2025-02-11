<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="icon" type="image/png" href="https://wishyicons.s3.eu-west-1.amazonaws.com/gift-box.png">
    <title>{{.ErrorCode}} - {{.ErrorMessage}}</title>
    <style>
        /* General Styles */
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background-color: #f8f9fa;
            text-align: center;
            padding: 50px;
            margin: 0;
        }

        /* Centered Box */
        .container {
            background-color: #fff;
            padding: 40px;
            max-width: 500px;
            margin: auto;
            border-radius: 12px;
            box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
            animation: fadeIn 0.5s ease-in-out;
        }

        /* Error Code */
        .error-code {
            font-size: 80px;
            font-weight: bold;
            color: #e74c3c;
            margin: 0;
        }

        /* Error Message */
        .error-message {
            font-size: 22px;
            color: #555;
            margin-bottom: 20px;
        }

        /* Return Button */
        .btn {
            display: inline-block;
            padding: 12px 20px;
            font-size: 18px;
            color: #fff;
            background: #3498db;
            text-decoration: none;
            border-radius: 8px;
            transition: 0.3s;
        }

        .btn:hover {
            background: #2980b9;
        }

        /* Fade-in Animation */
        @keyframes fadeIn {
            from {
                opacity: 0;
                transform: translateY(-10px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }

        /* Responsive Design */
        @media (max-width: 600px) {
            .container {
                padding: 30px;
            }

            .error-code {
                font-size: 60px;
            }

            .error-message {
                font-size: 18px;
            }
        }
    </style>
</head>
<body>

    <div class="container">
        <div class="error-code">{{.ErrorCode}}</div>
        <p class="error-message">Oops! {{.FriendlyErrorMessage}}</p>
    </div>

</body>
</html>
