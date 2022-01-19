
# README

This doc is the guideline for you to start up your fabric test network.
Please follow the steps below:

1. Copy the whole `fabric` dir to your own dir  such as `path/fabric`.
2. Exec `cd path/fabric && sh start.sh` to start up fabric network.
3. Register the deamon.sh into crontab to clear the chaincode regularly to invoide there is too mush data. 
   1. Edit  `crontab -e`
   2. Append `0 3 * * * cd path/fabric && sh deamon.sh > deamon.log`, which will clean fabric's data and then restart fabric network at 3:00 daily.
4. Edit the url in `fabric/config.yaml` to point to your network.

