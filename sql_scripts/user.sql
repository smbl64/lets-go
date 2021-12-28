CREATE user 'web'@'%';

GRANT SELECT, UPDATE, INSERT ON snippetbox.* to 'web'@'%';

ALTER USER 'web'@'%' IDENTIFIED BY 'web';
