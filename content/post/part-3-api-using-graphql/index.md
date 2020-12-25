---
title: Part 3:- API using GraphQL and Node.js
publishDate: 2018-12-31
draft: false
summary: This is a series of 3 articles which will help to write production grade code.
draft: false
authors:
- admin
tags:
- js
- opensource
- node.js
- docker
- mongodb
- testing
- grpahql
categories:
- opensource
- tech
image:
  placement: 2
  caption: "[Image](https://user-images.githubusercontent.com/24803604/50481874-edac8280-0a09-11e9-837b-3e64409c21fb.png)"
  focal_point: "Center"
  preview_only: false
  alt_text: API using GraphQL and Node.js
---

This is Article 3 for this series. You can find [first article here](https://knrt10.netlify.app/post/part-1-api-using-graphql/). and [second article here](https://knrt10.netlify.app/post/part-2-api-using-graphql/).

{{% toc %}}

## About this Article

In this article you will be able to do the following things:-

* Create an Auth based system.

* Use JWT to verify the tokens

* CRUD(Create Read Update Delete) operations for Todo using GraphQL

* Write tests for it.

## Creating a Auth based system

We need to have a system where users can create a todo only when they authenticate themselves as real user. For example, if we don’t authenticate our routes, anyone can create a todo and spammers can spam the database. Also we need to keep track of todo for a particular user. Including them and many other reasons we need to create an Auth based system. So let’s create a folder inside src called as middleware and then create 2 files inside it:-

* auth.ts

* index.ts

Inside `src/middleware/auth.ts` copy this :

```ts
import jwt = require("jsonwebtoken");
import { model } from "mongoose";
import { completeRequest } from "../functions/complete";
import {Response} from "../models";
import {UserSchema} from "../schemas";
import {Config} from "../shared";

const User = model("User", UserSchema);

/**
 * This is middleware to validate jwt token.
 * @param req
 * @param res
 * @param next
 */
export async function isAuthenticated(context) {
  const promise: Promise<Response> = new Promise((resolve, reject) => {
    const token: any = context.headers["x-access-token"];
    const secret: any = Config.secretKeys.jwtSecret;
    if (!token) {
      reject(new Response(403, "Auth token missing", {
        success: false,
      }));
    } else {
      // verify jwt token
      jwt.verify(token, secret, (err, decoded) => {
        if (err) {
          reject(new Response(500, "Unable to authenticate user", {
            success: false,
          }));
        } else {
          User.findById(decoded.id).select("password").then((user) => {
            if (!user) {
              reject(new Response(400, "Sorry No user found", {
                success: false,
              }));
            } else {
              resolve(new Response(200, "Successfull Response", {
                success: true,
                user,
                token,
              }));
            }
          });
        }
      });
    }
  });
  const val = await completeRequest(promise);
  return val;
}
```

**Explanation:-**

* *Line 1–6:-* We require necessary modules for this file.

* *Line 16:-* We are exporting a function to check whether a user is authenticated or not. This functions takes an argument which context which is generally a form of request.

* *Line 17:-* This function is returning a promise.

* *Line 18–23*:- Now as said earlier context is in form of request so it will contain headers, now we will check if any token is passed in x-access-token header. This is a custom header that we will use to pass our JWT token when making a GraphQL request to CRUD operation on a Todo. If token is not present in header we reject our promise with a newresponse that we created in response.ts.

* *Line 26–30:-* You can check [here](https://jwt.io/) how a JWT works. So if any modification is made to our token then it will reject the with a new response.

* *Line 32–48:-* If token is not tampered with then we recieve a decoded data. That decode data contains id of user. We then find that user with that id and return it. If no user is found we reject with a response.

* *Line 43–44:-* We resolve the promise by passing it to completeRequest and gets back another promise but this promise has response in form we want. We then return the value.

After 2 articles you might have idea what we are going to write in `src/middleware/index`.ts. You may have already written it 😄, but let’s be sure and copy this code.

```ts
export * from "./auth";
```

Now lets create our Schema for Todo for our mongoose model. Create a file inside schemas named `todoSchema.ts` and copy this inside

`src/schemas/todoSchema.ts`

```ts
import { Schema } from "mongoose";
import mongoose = require("mongoose");

mongoose.Promise = global.Promise;

/**
 * This is Schema for Blog
 * @constant {BlogSchema}
 */
export const TodoSchema = new Schema({
  id: {
    type: String,
  },
  postedByid: {
    type: String,
  },
  title: {
    type: String,
    required: true,
    select: true,
  },
  description: {
    type: String,
  },
  postedBy: {
    type: mongoose.Schema.Types.ObjectId,
    ref: "User",
  },
  name: {
    type: String, ref: "User",
  },
}, {
  timestamps: {},
});
```

We already understood how to make a schema when we created userSchema. Here is just one new thing on *line 25–28:-* We need to reference the todo that is created to the user who created it. So `postedBy` will contain `_id` of user who will create that todo. *Line 29–31* will store the name of the user.

Also update your `src/schemas/index.ts` and add

```ts
export * from "./todoSchema";
```

Now we will create functions for CRUD operations on a todo. So create file `todo.ts` inside routes folder and copy this inside

`src/routes/todo.ts`

```ts
import { model } from "mongoose";
import { completeRequest } from "../functions/complete";
import { Response } from "../models";
import { TodoSchema, UserSchema } from "../schemas";
const User = model("User", UserSchema);
const Todo = model("Todo", TodoSchema);

/**
 * This function creates a new Todo
 * @param args
 * @param user
 */
export async function addTodo(args, user) {
  const promise: Promise<Response> = new Promise((resolve, reject) => {
    if (!args.input.title || !args.input.description) {
      reject(new Response(200, "Please enter both title and description", {
        success: false,
      }));
    }

    const title = args.input.title.trim();
    const description = args.input.description.trim();
    if (!title.length || !description.length) {
      reject(new Response(200, "Title or description cannot be blank", {
        success: false,
      }));
    }

    User.findById({ _id: user._id }, (err, user) => {
      const todo = new Todo({
        postedBy: user.id,
        name: user.name,
        title,
        description,
      });

      todo.id = todo._id;
      todo.postedByid = user.id;
      todo.save((err) => {
        if (err) {
          reject(new Response(200, "Error in saving Todo", {
            success: false,
          }));
        }
        resolve(new Response(200, "Successfully saved Todo", {
          success: true,
          todo,
        }));
      });
    });
  });
  const val = await completeRequest(promise);
  return val;
}

/**
 * This returns all todos for user
 * @param user
 */
export async function getAlltodosForUser(user) {
  const promise: Promise<Response> = new Promise((resolve) => {
    Todo.find({ postedBy: user._id }, (err, todos) => {
      resolve(new Response(200, "All todos", {
        success: true,
        todos,
      }));
    });
  });
  const val = await completeRequest(promise);
  return val;
}

/**
 * This updates the todo information
 * @param user
 */
export async function update(args, user) {
  const promise: Promise<Response> = new Promise((resolve, reject) => {
    const todoId = args.input.id;
    if (!args.input.title || !args.input.description || !todoId) {
      reject(new Response(200, "Please enter all fields", {
        success: false,
      }));
    }

    const title = args.input.title.trim();
    const description = args.input.description.trim();
    if (!title.length || !description.length) {
      reject(new Response(200, "Title or description cannot be blank", {
        success: false,
      }));
    }

    Todo.findById({ _id: todoId }, (err, todo) => {
      if (err) {
        reject(new Response(200, "Not able to get todo", {
          success: false,
        }));
      } else if (String(todo.postedBy) !== String(user._id)) {
        reject(new Response(200, "You don't have access to update this todo", {
          success: false,
        }));
      } else {
        Todo.findOneAndUpdate({ _id: todoId }, { $set: { title, description } }, { new: true }, (err, todo) => {
          resolve(new Response(200, "Updated Todo", {
            success: true,
            todo,
          }));
        });
      }
    });
  });
  const val = await completeRequest(promise);
  return val;
}

/**
 * This function deletes the particular todo we want
 * @param args
 */
export async function deleteTodo(args, user) {
  const promise: Promise<Response> = new Promise((resolve, reject) => {
    const todoId = args.id;
    Todo.findById({ _id: todoId }, (err, data) => {
      if (err) {
        reject(new Response(200, "Not able to get todo", {
          success: false,
        }));
      } else if (!data) {
        reject(new Response(200, "Todo already deleted", {
          success: false,
        }));
      } else if (String(data.postedBy) !== String(user._id)) {
        reject(new Response(200, "You don't have access to delete this todo", {
          success: false,
        }));
      } else {
        Todo.findOneAndDelete({ _id: todoId }, () => {
          resolve(new Response(200, "Successfully deleted todo", {
            success: true,
          }));
        });
      }
    });
  });
  const val = await completeRequest(promise);
  return val;
}
```

You might now have idea what we are doing. We have 4 functions.

* *addTodo* :- This creates a new todo.

* *getAlltodosForUser* :- This returns all todos for an authenticated user.

* *update* :- This updates a particular Todo.

* *deleteTodo* :- This deletes a particular todo.

All Functions are authenticated. That is, only authenticated user can make request to perform CRUD operation.

Now lets create a graphQL schema for our Todo. Create a file `todographqlSchema.ts` inside schemas folder and copy this inside

`src/schemas/todographqlSchema.ts`

```ts
import { GraphQLBoolean, GraphQLID, GraphQLInputObjectType, GraphQLInt, GraphQLList, GraphQLNonNull, GraphQLObjectType, GraphQLString } from "graphql";

// User type
const todoType = new GraphQLObjectType({
  name: "todo",
  fields: {
    id: { type: GraphQLID },
    postedByid: { type: GraphQLID },
    name: { type: GraphQLString },
    title: { type: GraphQLString },
    description: { type: GraphQLString },
    createdAt: { type: GraphQLString },
    updatedAt: { type: GraphQLString },
  },
});

// Data reponse of user
const DataResponse = new GraphQLObjectType({
  name: "todoDataResponse",
  fields: {
    success: { type: GraphQLBoolean },
    todo: { type: todoType },
    token: { type: GraphQLString },
  },
});

// Response from User
const responseType = new GraphQLObjectType({
  name: "toDoResponse",
  fields: {
    code: { type: new GraphQLNonNull(GraphQLInt) },
    message: { type: new GraphQLNonNull(GraphQLString) },
    data: { type: new GraphQLNonNull(DataResponse) },
  },
});

const toDoInput = new GraphQLInputObjectType({
  name: "todoInput",
  fields: {
    title: { type: GraphQLString },
    description: { type: GraphQLString },
  },
});

const toDoInputUpdate = new GraphQLInputObjectType({
  name: "todoInputUpdate",
  fields: {
    id: { type: GraphQLString },
    title: { type: GraphQLString },
    description: { type: GraphQLString },
  },
});

// For getting todo

// Data reponse of user
const todoUsersDataResponse = new GraphQLObjectType({
  name: "todoUsersDataResponse",
  fields: {
    success: { type: GraphQLBoolean },
    todos: { type: new GraphQLList(todoType) },
  },
});

const userTodoResponse = new GraphQLObjectType({
  name: "userTodoResponse",
  fields: {
    code: { type: new GraphQLNonNull(GraphQLInt) },
    message: { type: new GraphQLNonNull(GraphQLString) },
    data: { type: new GraphQLNonNull(todoUsersDataResponse) },
  },
});

// For deleting Todo

const todoDeleteResponse = new GraphQLObjectType({
  name: "todoDeleteResponse",
  fields: {
    success: { type: GraphQLBoolean },
  },
});

const userTodoDeleteResponse = new GraphQLObjectType({
  name: "userTodoDeleteResponse",
  fields: {
    code: { type: new GraphQLNonNull(GraphQLInt) },
    message: { type: new GraphQLNonNull(GraphQLString) },
    data: { type: new GraphQLNonNull(todoDeleteResponse) },
  },
});

export const todographqlSchema = {
  todoType,
  DataResponse,
  responseType,
  toDoInput,
  toDoInputUpdate,
  todoUsersDataResponse,
  userTodoResponse,
  userTodoDeleteResponse,
};
```

You might have idea about it, as it is somewhat similar to `userLoginSchema.ts` and `userRegisterSchema.ts` for graphQL. Also we need to export this file so paste this line inside

`src/schemas/index.ts`

```ts
export * from "./todographqlSchema";
```

Now we need to finally add queries and mutation inside our `graphql.ts`. Like we did for `loginUser` and `registerUser`. So copy and replace your old code with this inside

`src/schemas/graphql.ts`

```ts
import { GraphQLNonNull, GraphQLObjectType, GraphQLSchema, GraphQLString } from "graphql";
import { isAuthenticated } from "../middleware";
import { addTodo, deleteTodo, getAlltodosForUser, login, register, update } from "../routes";
import { todographqlSchema } from "./todographqlSchema";
import { userLoginSchema } from "./userLoginSchema";
import { userRegisterSchema } from "./userRegisterSchema";

// Define the Query type
const queryType = new GraphQLObjectType({
  name: "Query",
  fields: {
    loginUser: {
      type: new GraphQLNonNull(userRegisterSchema.responseType),
      // `args` describes the arguments that the `user` query accepts
      args: {
        input: { type: userLoginSchema.UserInput },
      },
      async resolve(_, args) {
        const val = await login(args);
        return val;
      },
    },
    profileUser: {
      type: new GraphQLNonNull(userRegisterSchema.responseType),
      // `args` describes the arguments that the `user` query accepts
      async resolve(parent, args, context, info) {
        const authenticated = await isAuthenticated(context);
        return authenticated;
      },
    },
    todoUsers: {
      type: new GraphQLNonNull(todographqlSchema.userTodoResponse),
      // `args` describes the arguments that the `user` query accepts
      async resolve(parent, args, context, info) {
        const authenticated = await isAuthenticated(context);
        if (authenticated.code !== 200) {
          return authenticated;
        } else {
          const val = await getAlltodosForUser(authenticated.data.user);
          return val;
        }
      },
    },
  },
});

// Defining Mutation
const mutationType = new GraphQLObjectType({
  name: "Mutation",
  fields: {
    registerUser: {
      type: new GraphQLNonNull(userRegisterSchema.responseType),
      // `args` describes the arguments that the `user` query accepts
      args: {
        input: { type: userRegisterSchema.UserInput },
      },
      async resolve(_, args) {
        const val = await register(args);
        return val;
      },
    },
    addTodo: {
      type: new GraphQLNonNull(todographqlSchema.responseType),
      // `args` describes the arguments that the `user` query accepts
      args: {
        input: { type: todographqlSchema.toDoInput },
      },
      async resolve(parent, args, context, info) {
        const authenticated = await isAuthenticated(context);
        if (authenticated.code !== 200) {
          return authenticated;
        } else {
          const val = await addTodo(args, authenticated.data.user);
          return val;
        }
      },
    },
    updateTodo: {
      type: new GraphQLNonNull(todographqlSchema.responseType),
      // `args` describes the arguments that the `user` query accepts
      args: {
        input: { type: todographqlSchema.toDoInputUpdate },
      },
      async resolve(parent, args, context, info) {
        const authenticated = await isAuthenticated(context);
        if (authenticated.code !== 200) {
          return authenticated;
        } else {
          const val = await update(args, authenticated.data.user);
          return val;
        }
      },
    },
    deleteTodo: {
      type: new GraphQLNonNull(todographqlSchema.userTodoDeleteResponse),
      // `args` describes the arguments that the `user` query accepts
      args: {
        id: { type: GraphQLString },
      },
      async resolve(parent, args, context, info) {
        const authenticated = await isAuthenticated(context);
        if (authenticated.code !== 200) {
          return authenticated;
        } else {
          const val = await deleteTodo(args, authenticated.data.user);
          return val;
        }
      },
    },
  },
});

export const schema = new GraphQLSchema({
  query: queryType,
  mutation: mutationType,
});
```

Now lets write tests for our code changes. First we need to create 2 new files inside our test folder.

* profileQueries.ts

* todoQueries.ts

Inside `src/test/profileQueries.ts` copy this code

```ts
const query = `query profileUser{
  profileUser {
    code
    message
    data {
      success
      token
      user {
        name
      }
    }
  }
}`;

const profileSuccessfullyQuery = {
  query: query,
  operationName: "profileUser"
  ,
};

export const profileUser = {
  profileSuccessfullyQuery,
};
```

Also copy this to `src/test/todoQueries.ts`

```ts
const UserRoute = require("./user-test.spec");
const todoId = UserRoute.toDoSavedData;

const query = `mutation addTodo($input: todoInput) {
  addTodo(input: $input) {
    code
    message
    data {
      success
     	todo {
        id
        postedByid
        description
        updatedAt
        createdAt
        name
      }
    }
  }
}`;

const updateQuery = `mutation updateTodo($input: todoInputUpdate) {
  updateTodo(input: $input) {
    code
    message
    data {
      success
     	todo {
        id
        postedByid
        description
        updatedAt
        createdAt
        name
        title
      }
    }
  }
}`;

const deleteQuery = `mutation deleteTodo($id: String) {
  deleteTodo(id: $id) {
    code
    message
    data {
      success
    }
  }
}`;

const allTodos = `query todoUsers {
  todoUsers{
    code
    message
    data {
      success
      todos {
        title
        description
      }
    }
  }
}`;

const toDoSuccessfullyQuery = {
  query: query,
  operationName: "addTodo"
  ,
  variables: {
    input: {
      title: "Test title",
      description: "test description",
    },
  },
};

const toDoFailNotitleOrDescyQuery = {
  query: query,
  operationName: "addTodo"
  ,
  variables: {
    input: {
      title: "",
      description: "",
    },
  },
};

const toDoFailNotitleQuery = {
  query: query,
  operationName: "addTodo"
  ,
  variables: {
    input: {
      title: "   ",
      description: "dasdasda",
    },
  },
};

const toDoUpdateQuery = {
  query: updateQuery,
  operationName: "updateTodo"
  ,
  variables: {
    input: {
      id: "anything",
      title: "Test title",
      description: "test description",
    },
  },
};

const toDoFailNoIdUpdateQuery = {
  query: updateQuery,
  operationName: "updateTodo"
  ,
  variables: {
    input: {
      id: "",
      title: "fsdfs",
      description: "fsdsdf",
    },
  },
};

const toDoFailNotitleOrDescUpdateQuery = {
  query: updateQuery,
  operationName: "updateTodo"
  ,
  variables: {
    input: {
      id: "anything",
      title: "  ",
      description: "   ",
    },
  },
};

const toDoFailDeleteQuery = {
  query: deleteQuery,
  operationName: "deleteTodo"
  ,
  variables: {
    id: "anything",
  },
};

const TodoAllQuery = {
  query: allTodos,
  operationName: "todoUsers"
  ,
};

export const todoQueries = {
  toDoSuccessfullyQuery,
  toDoFailNotitleOrDescyQuery,
  toDoFailNotitleQuery,
  toDoUpdateQuery,
  toDoFailNoIdUpdateQuery,
  toDoFailNotitleOrDescUpdateQuery,
  toDoFailDeleteQuery,
  updateQuery,
  deleteQuery,
  TodoAllQuery,
};
view raw

```

We need to update our user-test.spec.ts file. This file is very big, so open the link given below and copy the file into your *[src/test/user-test.spec.ts](https://github.com/knrt10/Todo-backendAPI/blob/master/src/test/user-test.spec.ts)*

Lets run our test by running npm run build && npm run coverage. Make sure your MongoDB is up and running. You will get this output.

![Still 💯 code coverage 😉](https://cdn-images-1.medium.com/max/2880/1*xnnQGSOuzh5amiyAUu4SUw.png)

> *Still 💯 code coverage 😉*

## Accessing the API

First check your mongoDB server is up and running. Then start your server by running the following command

`npm start`

Now to access the API of application open your [GraphQL-Playground](https://github.com/prisma/graphql-playground) and enter url [http://localhost:3000/graphql](http://localhost:3000/graphql)

### Creating a Todo

Enter Query

```js
mutation addTodo($input: todoInput) {
  addTodo(input: $input) {
    code
    message
    data {
      success
      todo {
        id
        postedByid
        description
        updatedAt
        createdAt
        name
      }
    }
  }
}
```

and then query variable

```json
{
  "input": {
    "title": "This is a test todo",
    "description": "Lets see if this works"
  }
}
```

**Important:-** You need to set *HTTP HEADERS*. You can get this token with the use of login API.

```js
query loginUser($input: UserInputLogin) {
  loginUser(input: $input) {
    data {
      token
    }
  }
}
```

and then query variable

```json
{
  "input": {
    "username": "knrt10",
    "password": "test"
  }
}
```

From this response copy the token and then copy this code to your *HTTP HEADERS*

```json
{
  "x-access-token": "your access token from login API"
}
```

Then hit play button, you will get response like this. Try removing the token or changing token to something else and see the response.

![creating a Todo](https://cdn-images-1.medium.com/max/2880/1*1B0JXQ1Ht7_aua8XOuZKUw.png)*creating a Todo*

### Updating a Todo

Enter Query

```js
mutation updateTodo($input: todoInputUpdate) {
  updateTodo(input: $input) {
    code
    message
    data {
      success
      todo {
        id
        postedByid
        description
        updatedAt
        createdAt
        name
        title
      }
    }
  }
}
```

and then query variable. You can get the id from the todo you created before.

```json
{
  "input": {
    "id": "id of your todo you created before",
    "title": "Second check?",
    "description": "Yaaho."
  }
}
```

You also need to set header like above. After that when you hit play button you will get this response

![updating a Todo](https://cdn-images-1.medium.com/max/2880/1*Vx7r5F3dFriIgOchHbKIfg.png)*updating a Todo*

### Getting all Todo

Enter Query

```js
query todoUsers {
  todoUsers{
    code
    message
    data {
      success
      todos {
        title
        description
        id
      }
    }
  }
}
```

and then set *HTTP HEADERS* like above API.

```json
{
  "x-access-token": "Your access token"
}
```

After that when you hit play button you will get this response

![All todos for users](https://cdn-images-1.medium.com/max/2880/1*AadeoCY3Nbpp24uSCjgomQ.png)*All todos for users*

### Deleting a Todo

Enter Query

```js
mutation deleteTodo($id: String) {
  deleteTodo(id: $id) {
    code
    message
    data {
      success
    }
  }
}
```

and then query variable. You can get the id from the todo you created before.

```json
{
  "id": "5c25156f70d37365ede03609"
}
```

You also need to set header like above. After that when you hit play button you will get this response

![deleting a Todo](https://cdn-images-1.medium.com/max/2880/1*hqzW3K_Ie0M6Lp23J7B5NA.png)*deleting a Todo*

## Conclusion

That is for this part. In this part you learnt following things:-

* How to write Schema for GraphQL.

* How to create Schema for Todo and link it with a User

* Perform CRUD operation on a Todo using GraphQL.

* How to write clean code.

* How maintain 💯 code coverage 😉.

* You can see it took **33 tests** to run in about **804ms.** Which is less than a second. It shows how fast and precise our code is.

## **Code for this series**

Code is open sourced on [Github](https://github.com/knrt10/Todo-backendAPI) under MIT license. Feel free to use it as reference if you are stuck anywhere.

## Support

I wrote this series of articles by using my free time. A little motivation and support helps me a lot. If you like this nifty hack you can support me by doing any (or all 😉 ) of the following:

* Follow me on [Github](http://github.com/knrt10/) for more such projects.

* ⭐️ Star it on [Github](https://github.com/knrt10/Todo-backendAPI) and make it trend so that other people can know about my project.

### Did you find this page helpful? Consider sharing it 🙌
