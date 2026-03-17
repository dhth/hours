#!/bin/bash
sed -i 's/func DeleteActiveTL(db \*sql.DB) error {/func DeleteActiveTL(db \*sql.DB, beginTs time.Time) error {/' internal/persistence/queries.go
sed -i 's/WHERE active=true;/WHERE active=true AND begin_ts=?;/' internal/persistence/queries.go
sed -i 's/_, err = stmt.Exec()/_, err = stmt.Exec(beginTs.UTC())/' internal/persistence/queries.go

sed -i 's/func deleteActiveTL(db \*sql.DB) tea.Cmd {/func deleteActiveTL(db \*sql.DB, beginTs time.Time) tea.Cmd {/' internal/ui/cmds.go
sed -i 's/err := pers.DeleteActiveTL(db)/err := pers.DeleteActiveTL(db, beginTs)/' internal/ui/cmds.go

sed -i 's/cmds = append(cmds, deleteActiveTL(m.db))/cmds = append(cmds, deleteActiveTL(m.db, m.activeTLBeginTS))/' internal/ui/update.go
sed -i 's/return deleteActiveTL(m.db)/return deleteActiveTL(m.db, m.activeTLBeginTS)/' internal/ui/handle.go
