<!DOCTYPE html>
<html lang="en">

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
            color: #3498db;
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
            color: #2ecc71;
            margin-bottom: 15px;
            border-bottom: 2px solid #2ecc71;
            padding-bottom: 10px;
            text-align: center;
        }

        ul {
            list-style-type: none;
            padding: 0;
        }

        .card {
            background-color: #ffffff; /* Change the background color to a cheerful yellow */
            border: 1px solid #bdc3c7;
            margin: 15px 0;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
            display: flex;
            flex-wrap: wrap;
            justify-content: space-between;
        }

        .wish-item {
            width: 48%;
            margin-bottom: 10px;
        }

        .card a {
            color: #3498db;
            text-decoration: none;
            font-weight: bold;
            text-align: center;
            display: block;
            margin: 0 auto;
        }

        .card a:hover {
            text-decoration: underline;
        }
    </style>
</head>

<body>
    <h1>Niccolò Wish List</h1>

    {{range .}}
    <div class="category">
        <h2>{{.Cat}}</h2>
        <ul class="card">
            {{range .Wishes}}
            <li class="wish-item">
                <div class="card">
                    <a href="{{.Link}}" target="_blank">{{.Name}}</a>
                </div>
            </li>
            {{end}}
        </ul>
    </div>
    {{end}}

</body>

</html>
