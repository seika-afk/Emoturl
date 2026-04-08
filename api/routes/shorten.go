package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"emoturl/database"	


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
		
r2:= database.CreateClient(1);
defer r2.Close()
		val,err:=r2.Get(database.Ctx,c.IP()).Result()


if err == redis.Nil{
_ = r2.Set(database.Ctx,c.IP(),os.Get("API_QUOTA"),30*60*time.Second).Err()



}
else{


valInt,_ := strconv(val)

		if valInt<=0{


		limit,_:= r2.TTL(database.Ctx,c.IP()).Result()


			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{

				"error":"Rate Limit Exceeded",
				"rate_limit_reset":limit /time.Nanosecond


			})

		}



}




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




		r2.Dec(database.Ctx,c.IP())

	}







