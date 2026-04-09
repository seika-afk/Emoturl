package routes

import (
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"time"

	"emoturl/database"
	"emoturl/helpers"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL            string        `json:"url"`
	CustomShort    string        `json:"short"`
	Expiry         time.Duration `json:"expiry"`
	XRateRemaining int           `json:"rate_limit"`
	XRateLimitRest time.Duration `json:"rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error {
	body := new(request)

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "could not parse body",
		})
	}

	// implementing rate limiting
	r2 := database.CreateClient(1)
	defer r2.Close()

	val, err := r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":            "Rate Limit Exceeded",
				"rate_limit_reset": limit / time.Nanosecond,
			})
		}
	}

	// check if input is an actual url
	if _, err := url.ParseRequestURI(body.URL); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Incorrect URL",
		})
	}

	// check for domain server
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Domain Error",
		})
	}

	// enforce https,SSL
	body.URL = helpers.EnforceHTTP(body.URL)

	var shortCode string
	emotions := []string{"hahaha", "idcman", "What_am_I_even_doing_in_my_life", "fuck_this_shit", "Scam_link_hahahah", "dont_open_it"}
	if body.CustomShort == "" {
		shortCode = emotions[rand.Intn(len(emotions))]
	} else {
		shortCode = body.CustomShort
	}

	r := database.CreateClient(0)
	defer r.Close()

	val, _ = r.Get(database.Ctx, shortCode).Result()
	if val != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "URL custom short is already in use",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 24 * time.Hour
	}

	err = r.Set(database.Ctx, shortCode, body.URL, body.Expiry).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to connect to server",
		})
	}

	resp := response{
		URL:            body.URL,
		CustomShort:    "",
		Expiry:         body.Expiry,
		XRateRemaining: 10,
		XRateLimitRest: 30 * time.Minute,
	}

	_ = r2.Decr(database.Ctx, c.IP())

	val, _ = r2.Get(database.Ctx, c.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)

	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
	resp.XRateLimitRest = ttl

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + shortCode

	return c.Status(fiber.StatusOK).JSON(resp)
}
