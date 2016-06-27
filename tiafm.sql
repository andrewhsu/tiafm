-- tree.json
-- races-20150530.json
-- results-2900000000000071.json
CREATE TABLE IF NOT EXISTS race (
Season_Year INTEGER, -- 2015
Season_Id INTEGER, -- 9
Event_Code TEXT, -- 20150530/afm
Event_Name TEXT, -- May 30-31, 2015 @ Thunderhill
Event_Start_date TEXT, -- 2015-05-30
Event_End_date TEXT, -- 2015-05-31
Event_Track INTEGER, -- 2
Session_Combo_title TEXT, -- SUN R3
Session_Start TEXT, -- 2015-05-31 10:30:00
Session_Race_Actual_start TEXT, -- 2015-05-31 10:29:24
Session_Race_Actual_finish TEXT, -- 2015-05-31 10:43:24
Session_Race_ClassID TEXT, -- 115
Session_Race_ClassName TEXT, -- Nv 450 SB
Session_Race_Name TEXT, -- SUN R3 (Nv 450 SB)
Session_Race_Raid INTEGER, -- 2900000000000071
Session_Race_ResultStatus INTEGER, -- 2
Session_Race_Run_length INTEGER, -- 839917
Session_Race_Run_length_str INTEGER, -- 00:14:00
Session_Race_Sanction_id TEXT, -- 115
Session_Race_Scheduled_start TEXT, -- 2015-05-31 10:30:00
Session_Race_Scheduled_finish TEXT, -- 2015-05-31 10:50:00
Session_Race_Shortcode TEXT, -- YbWjIg
Session_Race_Status INTEGER, -- 2
Session_Race_Type TEXT, -- Race
Session_Race_Wave INTEGER, -- null
PRIMARY KEY (Session_Race_Raid)
);
CREATE TABLE IF NOT EXISTS result (
Race_Result_Ra_type TEXT, -- Race
Result_BestLap TEXT, -- 0:00.000
Result_BestLapOnLap INTEGER, -- 1
Result_ChampionshipPoints TEXT, -- 107
Result_Gap TEXT, -- 00.000
Result_LeaderGap TEXT, -- 5L
Result_Points INTEGER, -- 0
Result_Position INTEGER, -- 9994
Result_RacerID TEXT, -- 2500000000012086
Result_RacerName TEXT, -- Andrew Hsu
Result_Raid INTEGER, -- 2900000000000071
Result_Rrid TEXT, -- 2500000000017976
Result_SanctionID TEXT, -- 997
Result_SanctionStatus TEXT, -- N
Result_Sponsors TEXT, --
Result_Tiid TEXT, -- 2500000000012086
Result_TotalLaps INTEGER, -- 0
Result_TotalTime TEXT, -- 00.000
Result_VehicleID TEXT, -- 2500000000006942
Result_VehicleName TEXT, -- Kawasaki EX300
PRIMARY KEY (Result_Raid,Result_Tiid)
);
