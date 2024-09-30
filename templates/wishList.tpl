<!DOCTYPE html>
<html lang="it">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Niccolò Wish List</title>
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
            color: #111111; /* color of the bullet point (red in this example) */
            font-weight: bold; /* make the bullet point bold */
            display: inline-block;
            width: 1em; /* adjust spacing between bullet point and wish name */
            margin-left: 0.5em; /* align bullet point with previous headline */
        }
    </style>
</head>

<body>
    <h1>Niccolò Wish List</h1>

    {{range .}}
    <div class="category">
        <h2>{{.Cat}}</h2>
        <ul>
            {{range .Wishes}}
            <li class="wish-item">
                <a href="{{.Link}}" target="_blank">{{.Name}}</a>
            </li>
            {{end}}
        </ul>
    </div>
    {{end}}

</body>

</html>
