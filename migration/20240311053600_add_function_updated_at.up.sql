CREATE OR REPLACE FUNCTION trigger_set_updated()
  RETURNS trigger
  LANGUAGE plpgsql
AS $function$
BEGIN
	NEW.updatedAt = NOW();
	RETURN NEW;
END;
$function$