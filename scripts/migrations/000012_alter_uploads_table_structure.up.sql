ALTER TABLE uploads
    CHANGE COLUMN original_filename nama_file VARCHAR(255) NOT NULL,
    ADD COLUMN nama_file_sistem VARCHAR(255) NOT NULL DEFAULT '' AFTER nama_file,
    CHANGE COLUMN mime_type type_file VARCHAR(50) NOT NULL;
