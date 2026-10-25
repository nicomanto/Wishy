<!DOCTYPE html>
<html lang="it">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Gestione wishes · Wishy</title>
    <style>
        :root { color-scheme: dark; font-family: Inter, system-ui, sans-serif; background: #0f172a; color: #f1f5f9; }
        * { box-sizing: border-box; }
        body { margin: 0; padding: 32px 16px; }
        main { max-width: 1000px; margin: auto; }
        h1, h2 { margin-top: 0; }
        h1 { color: #c7d2fe; }
        h2 { font-size: 1.2rem; }
        section { margin: 20px 0; padding: 24px; border: 1px solid #334155; border-radius: 16px; background: #1e293b; }
        .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 16px; }
        .item { padding: 16px 0; border-top: 1px solid #334155; }
        label { display: block; margin: 12px 0 6px; color: #cbd5e1; }
        input, select, button { width: 100%; padding: 10px 12px; border: 1px solid #475569; border-radius: 8px; background: #0f172a; color: #f1f5f9; font: inherit; }
        input[type=checkbox] { width: auto; margin-right: 8px; }
        button { margin-top: 12px; border: 0; background: #818cf8; color: #0f172a; font-weight: 700; cursor: pointer; }
        button.danger { background: #fb7185; }
        .inline { display: inline; }
        .inline button { width: auto; }
        .error { padding: 12px; border-radius: 8px; background: #7f1d1d; }
        .muted { color: #94a3b8; }
        a { color: #a5b4fc; }
        @media (max-width: 600px) { body { padding: 16px 10px; } section { padding: 18px; } }
    </style>
</head>
<body>
<main>
    {{if .Authenticated}}
    <h1>Gestione wishes</h1>
    <p class="muted">Accesso effettuato come {{.User.Username}}.</p>
    <form class="inline" method="post" action="{{.Path}}">
        <input type="hidden" name="action" value="logout">
        <button type="submit">Esci</button>
    </form>
    {{if .Error}}<p class="error">{{.Error}}</p>{{end}}

    <section>
        <h2>Categorie</h2>
        <form method="post" action="{{$.Path}}">
            <input type="hidden" name="action" value="category_create">
            <label for="new-category">Nuova categoria</label>
            <input id="new-category" name="name" required maxlength="100">
            <button type="submit">Crea categoria</button>
        </form>
        {{range .Categories}}
        <div class="item">
            <form class="grid" method="post" action="{{$.Path}}">
                <input type="hidden" name="action" value="category_update">
                <input type="hidden" name="category_id" value="{{.ID.Hex}}">
                <label>Nome categoria<input name="name" value="{{.Name}}" required maxlength="100"></label>
                <button type="submit">Salva</button>
            </form>
            <form method="post" action="{{$.Path}}">
                <input type="hidden" name="action" value="category_delete">
                <input type="hidden" name="category_id" value="{{.ID.Hex}}">
                <button class="danger" type="submit">Elimina categoria</button>
            </form>
        </div>
        {{else}}<p class="muted">Non ci sono categorie.</p>{{end}}
    </section>

    <section>
        <h2>Wishes</h2>
        <form method="post" action="{{$.Path}}">
            <input type="hidden" name="action" value="wish_create">
            <div class="grid">
                <label>Nome<input name="name" required maxlength="200"></label>
                <label>Link<input name="link" type="url" maxlength="2048"></label>
                <label>Categoria
                    <select name="category_id" required>
                        {{range .Categories}}<option value="{{.ID.Hex}}">{{.Name}}</option>{{end}}
                    </select>
                </label>
                <label>Preferenza
                    <select name="preference">
                        <option value="1">Bassa</option><option value="2" selected>Media</option><option value="3">Alta</option>
                    </select>
                </label>
            </div>
            <label><input type="checkbox" name="active" checked> Wish attivo</label>
            <button type="submit">Crea wish</button>
        </form>
        {{range .Wishes}}
        {{$wish := .}}
        <div class="item">
            <form method="post" action="{{$.Path}}">
                <input type="hidden" name="action" value="wish_update">
                <input type="hidden" name="wish_id" value="{{.ID.Hex}}">
                <div class="grid">
                    <label>Nome<input name="name" value="{{.Name}}" required maxlength="200"></label>
                    <label>Link<input name="link" type="url" value="{{.Link}}" maxlength="2048"></label>
                    <label>Categoria
                        <select name="category_id" required>
                            {{range $.Categories}}<option value="{{.ID.Hex}}" {{if eq .ID $wish.Category.ID}}selected{{end}}>{{.Name}}</option>{{end}}
                        </select>
                    </label>
                    <label>Preferenza
                        <select name="preference">
                            <option value="1" {{if eq .Preference 1}}selected{{end}}>Bassa</option>
                            <option value="2" {{if eq .Preference 2}}selected{{end}}>Media</option>
                            <option value="3" {{if eq .Preference 3}}selected{{end}}>Alta</option>
                        </select>
                    </label>
                </div>
                <label><input type="checkbox" name="active" {{if .Active}}checked{{end}}> Wish attivo</label>
                <button type="submit">Salva wish</button>
            </form>
            <form method="post" action="{{$.Path}}">
                <input type="hidden" name="action" value="wish_delete">
                <input type="hidden" name="wish_id" value="{{.ID.Hex}}">
                <button class="danger" type="submit">Elimina wish</button>
            </form>
        </div>
        {{else}}<p class="muted">Non ci sono wishes.</p>{{end}}
    </section>
    {{else}}
    <section>
        <h1>Accedi a Wishy</h1>
        {{if .Error}}<p class="error">{{.Error}}</p>{{end}}
        <form method="post" action="{{.Path}}">
            <input type="hidden" name="action" value="login">
            <label for="username">Username</label>
            <input id="username" name="username" autocomplete="username" required>
            <label for="password">Password</label>
            <input id="password" name="password" type="password" autocomplete="current-password" required>
            <button type="submit">Accedi</button>
        </form>
    </section>
    {{end}}
</main>
</body>
</html>
