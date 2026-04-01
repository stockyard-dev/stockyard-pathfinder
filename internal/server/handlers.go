package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-pathfinder/internal/store")
func(s *Server)handleList(w http.ResponseWriter,r *http.Request){list,_:=s.db.List();if list==nil{list=[]store.Rule{}};writeJSON(w,200,list)}
func(s *Server)handleCreate(w http.ResponseWriter,r *http.Request){var rule store.Rule;json.NewDecoder(r.Body).Decode(&rule);if rule.Pattern==""||rule.Target==""{writeError(w,400,"pattern and target required");return};if rule.Type==""{rule.Type="redirect"};if rule.StatusCode==0{rule.StatusCode=301};rule.Active=true;s.db.Create(&rule);writeJSON(w,201,rule)}
func(s *Server)handleToggle(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.Toggle(id);writeJSON(w,200,map[string]string{"status":"toggled"})}
func(s *Server)handleDelete(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.Delete(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
