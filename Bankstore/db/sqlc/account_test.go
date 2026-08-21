package db

/*
	Не учтены в проверках многие кейсы (в рамках обучения)
	Ex. проверить, что при обновлении не меняется ничего кроме баланса и т.п.
*/

import (
	"context"
	"testing"

	"Bankstore/utils"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

// Создаём аккаунт, ищём его же, должны быть идентичны
func TestGetAccount(t *testing.T) {
	createdAccount := createRandomAccount(t)
	foundAccount, err := testQueries.GetAccount(context.Background(), createdAccount.ID)

	require.NoError(t, err)
	require.Equal(t, createdAccount, foundAccount)
}

// Создадим 5 аккаунтов, запросим список из трёх. Убедимся, что вернулось три не пустых объекта
// Оффсет в рамках данной задачи не будем проверять, т.к. не знаем состояние данных в БД на момент запуска.
// Теоретически можно взять список всех аккаунтов и его размер использовать для проверки оффсета (при большом оффсете должно возвращаться данных меньше лимита)
func TestListAccounts(t *testing.T) {
	for range 5 {
		createRandomAccount(t)
	}

	accounts, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{
		Limit:  3,
		Offset: 0,
	})
	require.NoError(t, err)
	require.Len(t, accounts, 3)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

// Создаем аккаунт с фиксированными данными, обновляем данные, получаем ещё раз из БД и убеждаемся, что данные обновились
func TestUpdateAccount(t *testing.T) {
	user := createRandomUser(t)
	createdAccount, err := testQueries.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    user.Username,
		Balance:  1000,
		Currency: "USD",
	})
	require.NoError(t, err)
	require.EqualValues(t, createdAccount.Balance, 1000)

	updatedAccount, err := testQueries.UpdateAccount(context.Background(), UpdateAccountParams{
		ID:      createdAccount.ID,
		Balance: 5000,
	})
	require.NoError(t, err)
	require.EqualValues(t, updatedAccount.Balance, 5000)

	accountFromDB, err := testQueries.GetAccount(context.Background(), createdAccount.ID)
	require.NoError(t, err)
	require.EqualValues(t, accountFromDB.Balance, 5000)
}

// Создаем аккаунт и удаляем его. Пытаемся найти его в БД.
func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	cnt, err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.EqualValues(t, 1, cnt)

	notFoundAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.Empty(t, notFoundAccount)
}

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	ra := utils.RandomAccount()
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  ra.Balance,
		Currency: Currency(ra.Currency),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}
