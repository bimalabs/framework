# Model Relationship in Bima

Bima using Gorm to handle database connectivity, and also use gRPC Gateway to convert gRPC to JSON response.

In Gorm, for relationship, column name basically named with `xxxID` but for gRPC, column named using camel case standard.

We decide to combine the stadards and use `xxxId` as relationship indicator in Gorm, and then the json will map to `xxxId` in response.
