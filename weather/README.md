Информация о названиях городов и их координаты записывается в базу данных автоматически при запуске сервиса (в main.go файле находится
[]cityList).

В файле .env содержится необходимая конфигурация для запуска сервиса (может читаться из переменных окружения), в том числе и API key  для работы со сторонним сервисом


EndPoints list for User:

"/weather/locations" -  list of cities, which we can get forcast

"/weather/forecast/{IdOfLocation}" - forecast for city with dates we can get

"/weather/locations/{IdOfLocation}?time=" - detailForecast by time      // need to use By time and locationID