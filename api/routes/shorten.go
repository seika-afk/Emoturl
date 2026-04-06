package routes

import (
	"time"

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


func ShortenURL(c *fiber.Ctx) error{

		body := new(request)

		if err := c.BodyParser(body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "could not parse body",
			})

		//implementing rate limiting 
		//check if input is an actual url 

		if !govalidator.IsURL(body.URL){
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Incorrect URL",
		}


		// check for domain server 
		if !helpers.RemoveDomainError(body.URL){

			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Domain Error",
	


		}


		// enforce https,SSL
		body.URL = helpers.EnforceHTTP(body.URL)

	}







