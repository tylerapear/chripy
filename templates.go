package main


func loadTemplates() map[string]string {
   
    return map[string]string{
        "admin_metrics_template": `

<html>
    <body>
        <h1>Welcome, Chirpy Admin</h1>
        <p>Chirpy has been visited %d times!</p>
    </body>
</html>

`,

    }
}
