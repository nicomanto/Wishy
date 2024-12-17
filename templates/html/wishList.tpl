<!DOCTYPE html>
<html lang="it">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="icon" type="image/png" href="https://wishyicons.s3.eu-west-1.amazonaws.com/gift-box.png">
    <title>{{.Username}} Wish List</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            margin: 20px;
            padding: 20px;
            background-color: #f9f9f9;
            color: #333;
        }

        h1 {
            color: #2e3330;
            text-align: center;
            margin-bottom: 20px;
        }

        .last-update {
            text-align: center;
            font-size: 0.8em;
            color: #666;
            margin-top: 0;
            margin-bottom: 20px;
        }

        .category {
            background-color: #ffffff;
            margin: 20px 0;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        }

        h2 {
            color: #666766;
            margin-bottom: 15px;
            border-bottom: 2px solid #666766;
            padding-bottom: 10px;
            text-align: center;
        }

        ul {
            list-style-type: none;
            padding: 0;
        }

        .wish-item {
            margin-bottom: 10px;
        }

        .wish-item a {
            color: #3498db;
            text-decoration: none;
            font-weight: bold;
        }

        .wish-item a:hover {
            text-decoration: underline;
        }

        .wish-item::before {
            content: "\2022"; /* bullet point character */
            color: #666766; /* color of the bullet point */
            font-weight: bold; /* make the bullet point bold */
            display: inline-block;
            width: 1em; /* adjust spacing between bullet point and wish name */
            margin-left: 0.5em; /* align bullet point with previous headline */
        }
        
        .preference {
            font-weight: bold;
            border-radius: 5px;
            padding: 3px 8px;
            font-size: 0.9em;
            color: #fff;
        }
    </style>
</head>

<body>
    <h1>{{.Username}} Wish List</h1>
    <div class="last-update">Last update: {{.LastUpdate}}</div>

    {{range .Wishes}}
    <div class="category">
        <h2>{{.Cat}}</h2>
        <ul>
            {{range .Wishes}}
            <li class="wish-item">
                <a href="{{.Link}}" target="_blank">{{.Name}}</a>
                {{if eq .Preference 3}}
                    <span class="preference">&#x2B50;&#x2B50;&#x2B50;</span>
                {{end}}
                {{if eq .Preference 2}}
                    <span class="preference">&#x2B50;&#x2B50</span>
                {{end}}
                {{if eq .Preference 1}}
                    <span class="preference">&#x2B50;</span>
                {{end}}
            </li>
            {{end}}
        </ul>
    </div>
    {{end}}

</body>

</html>
