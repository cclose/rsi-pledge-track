CREATE TABLE IF NOT EXISTS PledgeData (
    `ID` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `TimeStamp` datetime NOT NULL,
    `Funding` int(11) unsigned DEFAULT NULL,
    `Citizens` int(11) unsigned DEFAULT NULL,
    `Fleet` int(11) unsigned DEFAULT NULL,
    PRIMARY KEY (`ID`),
    UNIQUE KEY `TimeStamp_UNIQUE` (`TimeStamp`)
) ENGINE=InnoDB CHARACTER SET=utf8mb4 COLLATE=utf8mb4_general_ci;