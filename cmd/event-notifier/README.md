Events collection needs to be created with capped settings

```shell
# connect to the mongosh
$ mongosh "mongodb://username:password@host:27017/"
test> use pac
pac> db.createCollection("events", { capped: true, size: 300000 })
{ ok: 1 }
pac>
```