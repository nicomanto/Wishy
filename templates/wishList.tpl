<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Niccolò's Wish List</title>
    <style>
        body {
            font-family: 'Arial', sans-serif;
            margin: 20px;
            padding: 20px;
            background-color: #f4f4f4;
        }

        h1 {
            color: #333;
        }

        ul {
            list-style-type: none;
            padding: 0;
        }

        li {
            background-color: #fff;
            border: 1px solid #ddd;
            margin: 10px 0;
            padding: 15px;
            border-radius: 5px;
            box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
        }

        strong {
            color: #333;
        }

        a {
            color: #428bca;
            text-decoration: none;
            font-weight: bold;
        }

        a:hover {
            text-decoration: underline;
        }
    </style>
</head>

<body>
    <h1>Niccolò's Wish List</h1>
    <ul>
        {{range .}}
        <li>
            <strong>Name:</strong> {{.Name}}<br>
            <strong>Link:</strong> <a href="{{.Link}}" target="_blank">{{.Link}}</a><br>
            <strong>Category:</strong> {{.Category.Name}}<br>
        </li>
        {{end}}
    </ul>
</body>

</html>
