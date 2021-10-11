conn = new Mongo();
db = conn.getDB("resort");

db = connect("localhost:27017/resort");

db.users.drop();

db.createCollection("users", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["firstName", "lastName", "email", "password", "phoneNumber"],
            properties: {
                hotel: {
                    bsonType: "array",
                    description: "reserved hotel rooms"
                },
                restaurant: {
                    bsonType: "array",
                    description: "ordered foods"
                },
                profile: {
                    firstName: {
                        bsonType: "string",
                        description: "must be a string and is required"
                    },
                    lastName: {
                        bsonType: "string",
                        description: "must be a string and is required"
                    },
                    email: {
                        unique: true,
                        bsonType: "string",
                        // Regexp to validate emails with more strict rules as added in tests/users.js which also conforms mostly with RFC2822 guide lines
                        match: [/^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/, 'Please enter a valid email'],
                    },
                    password: {
                        bsonType: "string",
                        description: "must be a string and is required"
                    },
                    phoneNumber: {
                        bsonType: "string",
                        description: "must be a string and is required"
                    },
                    bsonType: "object",
                    description: "more personal information"
                },

            }
        }
    }
});

db.users.insertOne({profile: { firstName: "arian", lastName: "pourarian", email: "a.pourarian@gm.com", password: "13731892", phoneNumber: "00989054778974" }});
db.users.insertOne({profile: { firstName: "sahel", lastName: "shamsi", email: "s.shamsi@yh.com", password: "12345678", phoneNumber: "00989356231225" }});