CREATE TABLE userdata_public
(
  path text NOT NULL, -- 'site/username/paramName': 'freefeed.net/davidmz/selfdescr'
  value json NOT NULL,
  CONSTRAINT userdata_public_pkey PRIMARY KEY (path)
)
WITH (
  OIDS=FALSE
);
