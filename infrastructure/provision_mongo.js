use manageme
db.createUser({user:"mdbmanageme", pwd:"manageme", roles:["readWrite"]})

use test
db.createUser({user:"mdbmanageme", pwd:"manageme", roles:["readWrite"]})
