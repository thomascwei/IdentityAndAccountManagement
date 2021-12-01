package token

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker("12345678123456781234567812345678")
	require.NoError(t, err)

	username := "thomas"
	duration := time.Second

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	time.Sleep(duration * 2)
	token, err := maker.CreateToken(username, duration)
	fmt.Println(token)
	//v2.local.MVU_4-QBa0Dx29WY5cvmejYieOzbyC-kn1-_XToSpVGH_CppRYVdawTKfc89rJqfnKFlkShF__uE0orjFOM9als5QKgS_4IPaLrmoangaj_mltaTWj2O_XP8F3vfdKyxDno2ANpESO0Ga9SgnSQXOhnCWP3ydmGb4T3SXgnQ9NvqWMHOveYtiGAl-rjYxAXwliiwoDNZW5omGJYYcgPKaKyql1UWYgOtxMqOm4VZ3XayNEZ-Tz-vPTTs7WhLyYYgL7nzTehgi0wADg.bnVsbA
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker("12345678123456781234567812345678")
	require.NoError(t, err)

	token, err := maker.CreateToken("thomas", -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
