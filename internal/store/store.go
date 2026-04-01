package store
import("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{*sql.DB}
type Rule struct{ID int64 `json:"id"`;Pattern string `json:"pattern"`;Target string `json:"target"`;Type string `json:"type"`;StatusCode int `json:"status_code"`;HitCount int64 `json:"hit_count"`;Active bool `json:"active"`;CreatedAt time.Time `json:"created_at"`}
func Open(d string)(*DB,error){os.MkdirAll(d,0755);dsn:=filepath.Join(d,"pathfinder.db")+"?_journal_mode=WAL&_busy_timeout=5000";db,err:=sql.Open("sqlite",dsn);if err!=nil{return nil,fmt.Errorf("open: %w",err)};db.SetMaxOpenConns(1);migrate(db);return &DB{db},nil}
func migrate(db *sql.DB){db.Exec(`CREATE TABLE IF NOT EXISTS rules(id INTEGER PRIMARY KEY AUTOINCREMENT,pattern TEXT NOT NULL,target TEXT NOT NULL,type TEXT DEFAULT 'redirect',status_code INTEGER DEFAULT 301,hit_count INTEGER DEFAULT 0,active INTEGER DEFAULT 1,created_at DATETIME DEFAULT CURRENT_TIMESTAMP)`)}
func(db *DB)Create(r *Rule)error{act:=1;if !r.Active{act=0};res,err:=db.Exec(`INSERT INTO rules(pattern,target,type,status_code,active)VALUES(?,?,?,?,?)`,r.Pattern,r.Target,r.Type,r.StatusCode,act);if err!=nil{return err};r.ID,_=res.LastInsertId();return nil}
func(db *DB)List()([]Rule,error){rows,_:=db.Query(`SELECT id,pattern,target,type,status_code,hit_count,active,created_at FROM rules ORDER BY created_at DESC`);defer rows.Close();var out[]Rule;for rows.Next(){var r Rule;var act int;rows.Scan(&r.ID,&r.Pattern,&r.Target,&r.Type,&r.StatusCode,&r.HitCount,&act,&r.CreatedAt);r.Active=act==1;out=append(out,r)};return out,nil}
func(db *DB)Toggle(id int64){db.Exec(`UPDATE rules SET active=1-active WHERE id=?`,id)}
func(db *DB)Hit(id int64){db.Exec(`UPDATE rules SET hit_count=hit_count+1 WHERE id=?`,id)}
func(db *DB)Delete(id int64){db.Exec(`DELETE FROM rules WHERE id=?`,id)}
func(db *DB)Stats()(map[string]interface{},error){var total,active int64;db.QueryRow(`SELECT COUNT(*),SUM(hit_count) FROM rules`).Scan(&total,&active);return map[string]interface{}{"rules":total,"total_hits":active},nil}
