# time-to-go
Command line tool to display upcoming NYC subway times on a Raspberry Pi OLED monitor

## Command Usage
Print next subway times to OLED display:
```
./time-to-go --mta-api-key YOUR_MTA_API_KEY --stop-id YOUR_STOP_ID [--bus YOUR_I2C_BUS]
```

Print next subway times to terminal:
```
./time-to-go --debug --mta-api-key YOUR_MTA_API_KEY --stop-id YOUR_STOP_ID [--bus YOUR_I2C_BUS]
```

Clear OLED display:
```
./time-to-go clear [--bus YOUR_I2C_BUS]
```

![IMG_2382 (1)](https://github.com/christian-doucette/time-to-go/assets/64502867/cd741130-a423-456a-b77c-278f551a23d2)

## Installing
1. Install Go on your Raspberry Pi, if you don't have it already
   
Follow [this guide](https://www.jeremymorgan.com/tutorials/raspberry-pi/install-go-raspberry-pi/) to install Go. I used ```go1.22.0.linux-arm64.tar.gz``` for my Raspberry Pi 4 running a 64 bit OS.

2. Download the repo to your Raspberry Pi
```
git clone https://github.com/christian-doucette/time-to-go.git
```
3. Build the executable

Navigate into the repo then run:
```
go build
```
This should create an executable ```time-to-go``` file

4. Request an API key from the MTA

Sign up for one [here](https://api.mta.info/#/signup).

5. Test run (print to output)

At this point, try running the following command, using your MTA api key:

```
./time-to-go --debug --mta-api-key XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX --stop-id R16N
```
This will print the next subway times at stop R16N (Times Square–42nd Street NQRW) to output. If something goes wrong, resolve it before continuing.



 
5. Connect the OLED display

I used a 0.96 inch display with an SSD1306 driver connecting over I2C. Follow the beginning part of [this guide](https://www.raspberrypi-spy.co.uk/2018/04/i2c-oled-display-module-with-raspberry-pi/) for wiring and setup of the OLED display. If you want to use multiple monitors, follow [this guide](https://www.instructables.com/Raspberry-PI-Multiple-I2c-Devices/) for adding additional I2C buses and wiring the additional displays. 


6. Test run (print to OLED display)

At this point, try running the following command:
```
./time-to-go --mta-api-key XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX --stop-id R16N
```

This should display the next subway times at stop R16N to your OLED display. If something goes wrong, resolve it before continuing.

You can run the following command to clear the display:
```
./time-to-go clear
```

7. Set up a Cronjob to run the command continuously

For a refreshing display, you will want to rerun the command every 30 seconds to ensure the OLED dispaly always has up to date info.

Edit cronjobs:
```
crontab -e
```

Add the following lines:
```
* * * * *            ~/Desktop/time-to-go/time-to-go --mta-api-key XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX --stop-id R16N
* * * * * (sleep 30; ~/Desktop/time-to-go/time-to-go --mta-api-key XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX --stop-id R16N)
```
Cronjobs can be run at most once per minute, so to run it every thirty seconds there are two Cronjobs kicked off on the minute, one which waits 30 seconds before execution.
- Replace the mta-api-key argument with your key
- Replace the stop-id argument with the ID of the stop you want to display (find a stop ID by name [here](https://github.com/christian-doucette/time-to-go/blob/main/internal/gtfs/stops.txt)).
- Replace ``` ~/Desktop/time-to-go/time-to-go``` with the full path to the executable

If you are using multiple I2C devices on different buses, specify the buses (1 will be used as the default):
```
* * * * *            ~/Desktop/time-to-go/time-to-go --mta-api-key XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX --stop-id R16N --bus 1
* * * * * (sleep 30; ~/Desktop/time-to-go/time-to-go --mta-api-key XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX --stop-id R16N --bus 1)
* * * * *            ~/Desktop/time-to-go/time-to-go --mta-api-key XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX --stop-id R16S --bus 3
* * * * * (sleep 30; ~/Desktop/time-to-go/time-to-go --mta-api-key XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX --stop-id R16S --bus 3)

```








