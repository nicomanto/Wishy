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

        /* Page Header */
        .container {
            max-width: 600px;
            margin: 40px auto;
            padding: 20px;
        }

        h1 {
            font-size: 28px;
            color: #2c3e50;
        }

        .last-update {
            font-size: 14px;
            color: #777;
            margin-bottom: 30px;
        }

        /* Category Card */
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

        /* Wishlist Items */
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

        /* Preference Stars */
        .preference {
            font-size: 18px;
            color: #f1c40f;
        }

        /* Responsive Design */
        @media (max-width: 600px) {
            .container {
                padding: 10px;
            }
        }
    </style>
</head>

<body>
    <div class="container">
        <h1>🎁 {{.Username}}'s Wish List</h1>
        <div class="last-update">Last updated: {{.LastUpdate}}</div>

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

</html>
