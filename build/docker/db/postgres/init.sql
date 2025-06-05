create table public.specialization
(
    id    integer not null
        constraint specialization_pk
            primary key,
    title varchar not null
);

alter table public.specialization
    owner to yydz;

create table public.cabinet
(
    id                uuid    not null
        constraint cabinet_pk
            primary key,
    number            integer not null,
    specialization_id integer not null
        constraint cabinet_specialization_id_fk
            references public.specialization
);

alter table public.cabinet
    owner to yydz;

create table public.doctor
(
    id                uuid    not null
        constraint doctor_pk
            primary key,
    last_name         varchar not null,
    first_name        varchar not null,
    second_name       varchar,
    avatar_link       varchar,
    specialization_id integer not null
        constraint doctor_specialization_id_fk
            references public.specialization
);

alter table public.doctor
    owner to yydz;

create table public.profile
(
    id                  uuid    not null
        constraint profile_pk
            primary key,
    phone               integer not null,
    last_name           varchar not null,
    first_name          varchar not null,
    second_name         varchar,
    gender              boolean not null,
    birth_date          date    not null,
    passport_series     varchar not null,
    passport_number     varchar not null,
    passport_issuer     varchar not null,
    passport_issue_date date    not null,
    snils               varchar not null
);

alter table public.profile
    owner to yydz;

create table public.appointment
(
    date       date    not null,
    time       varchar not null,
    cabinet_id uuid    not null
        constraint appointment_cabinet_id_fk
            references public.cabinet,
    doctor_id  uuid    not null
        constraint appointment_doctor_id_fk
            references public.doctor,
    profile_id uuid
        constraint appointment_profile_id_fk
            references public.profile,
    constraint appointment_pk
        primary key (time, date, doctor_id)
);

alter table public.appointment
    owner to yydz;