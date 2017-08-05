# london-audio-microservice
[![CircleCI](https://img.shields.io/circleci/project/byuoitav/london-audio-microservice.svg)](https://circleci.com/gh/byuoitav/london-audio-microservice) [![Apache 2 License](https://img.shields.io/hexpm/l/plug.svg)](https://raw.githubusercontent.com/byuoitav/london-audio-microservice/master/LICENSE)

## DSP Control
The microservice controls and monitors gain blocks inside Soundweb London devices over TCP. 
The address of the gain blocks must be passed in as a hex string that corresponds to the address to a gain block that already exists on the London DSP
