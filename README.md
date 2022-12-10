# mqtt client to support French electriciy suppliers contract named "Tempo"

## How to run it

few environment variables must be set:

- TEMPO_URL</br>
  should be filled with <https://particulier.edf.fr/services/rest/referentiel/searchTempoStore?dateRelevant=%s&TypeAlerte=TEMPO></br>
  where %s will be replaced by the current date (format 2006-01-02)
- TEMPO_MQTT_URL</br>
  address of the mqtt broker, looks like tcp://localhost:1883
- TEMPO_MQTT_LOGIN</br>
  login for the broker. "" if none.
- TEMPO_MQTT_PASSWORD</br>
  password for the broker. "" if none.
- DEBUG</br>
  set to anything but empty to get more messages

for instance:

    DEBUG=true TEMPO_URL='https://particulier.edf.fr/services/rest/referentiel/searchTempoStore?dateRelevant=%s&TypeAlerte=TEMPO' TEMPO_MQTT_URL="tcp://localhost:1883" TEMPO_MQTT_LOGIN="" TEMPO_MQTT_PASSWORD="" ./release/1.0.0/linux/amd64/mqtt-tempo

## How to install it

An ansible role is provided here for RHEL: <https://bitbucket.org/mbee/ansible-role-mqtt-tempo>
