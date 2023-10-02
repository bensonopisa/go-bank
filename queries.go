package main

const (
	createTable = `create table if not exists accounts(
		id serial not null primary key,
		name varchar(50),
		balance int,
		createdAt timestamp
	)`

	fetchAllAccounts = `select * from accounts`

	insertNewAccount = `insert into accounts (name, balance, createdAt) values ($1, $2, $3)`

	updateAccount = `update accounts set name = $1 where id =$2`
)