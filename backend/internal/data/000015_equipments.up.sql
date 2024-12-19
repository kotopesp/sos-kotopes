CREATE TABLE IF NOT EXISTS
    equipments
(
    id                  SERIAL PRIMARY KEY,
    have_metal_cage     BOOLEAN,
    have_plastic_cage   BOOLEAN,
    have_net            BOOLEAN,
    have_ladder         BOOLEAN,
    have_other          VARCHAR
    );

ALTER TABLE IF EXISTS seekers
    ADD equipment_id  INTEGER REFERENCES equipments (id);
