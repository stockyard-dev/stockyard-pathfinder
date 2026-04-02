package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Route struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Pattern string `json:"pattern"`
	Target string `json:"target"`
	Method string `json:"method"`
	Priority int `json:"priority"`
	RateLimit int `json:"rate_limit"`
	Enabled int `json:"enabled"`
	HitCount int `json:"hit_count"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"pathfinder.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS routes(id TEXT PRIMARY KEY,name TEXT NOT NULL,pattern TEXT DEFAULT '',target TEXT DEFAULT '',method TEXT DEFAULT '*',priority INTEGER DEFAULT 0,rate_limit INTEGER DEFAULT 0,enabled INTEGER DEFAULT 1,hit_count INTEGER DEFAULT 0,created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Route)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO routes(id,name,pattern,target,method,priority,rate_limit,enabled,hit_count,created_at)VALUES(?,?,?,?,?,?,?,?,?,?)`,e.ID,e.Name,e.Pattern,e.Target,e.Method,e.Priority,e.RateLimit,e.Enabled,e.HitCount,e.CreatedAt);return err}
func(d *DB)Get(id string)*Route{var e Route;if d.db.QueryRow(`SELECT id,name,pattern,target,method,priority,rate_limit,enabled,hit_count,created_at FROM routes WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.Pattern,&e.Target,&e.Method,&e.Priority,&e.RateLimit,&e.Enabled,&e.HitCount,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Route{rows,_:=d.db.Query(`SELECT id,name,pattern,target,method,priority,rate_limit,enabled,hit_count,created_at FROM routes ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Route;for rows.Next(){var e Route;rows.Scan(&e.ID,&e.Name,&e.Pattern,&e.Target,&e.Method,&e.Priority,&e.RateLimit,&e.Enabled,&e.HitCount,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Update(e *Route)error{_,err:=d.db.Exec(`UPDATE routes SET name=?,pattern=?,target=?,method=?,priority=?,rate_limit=?,enabled=?,hit_count=? WHERE id=?`,e.Name,e.Pattern,e.Target,e.Method,e.Priority,e.RateLimit,e.Enabled,e.HitCount,e.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM routes WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM routes`).Scan(&n);return n}

func(d *DB)Search(q string, filters map[string]string)[]Route{
    where:="1=1"
    args:=[]any{}
    if q!=""{
        where+=" AND (name LIKE ?)"
        args=append(args,"%"+q+"%");
    }
    if v,ok:=filters["priority"];ok&&v!=""{where+=" AND priority=?";args=append(args,v)}
    if v,ok:=filters["enabled"];ok&&v!=""{where+=" AND enabled=?";args=append(args,v)}
    rows,_:=d.db.Query(`SELECT id,name,pattern,target,method,priority,rate_limit,enabled,hit_count,created_at FROM routes WHERE `+where+` ORDER BY created_at DESC`,args...)
    if rows==nil{return nil};defer rows.Close()
    var o []Route;for rows.Next(){var e Route;rows.Scan(&e.ID,&e.Name,&e.Pattern,&e.Target,&e.Method,&e.Priority,&e.RateLimit,&e.Enabled,&e.HitCount,&e.CreatedAt);o=append(o,e)};return o
}

func(d *DB)Stats()map[string]any{
    m:=map[string]any{"total":d.Count()}
    return m
}
