package task_campaign

import (
    "context"

    "github.com/pkg/errors"

    "safeweb.app/config"
    "safeweb.app/service/safeweb_lib/utils"
)

func FlattenBaseInfo(ctx context.Context, input interface{}, cfg interface{}) error {
    if service, err := load(cfg); err != nil {
        return err
    } else {
        if slices, err := safeweb_lib_utils.MakeSliceFromInput(input); err != nil {
            return err
        } else {
            for _, value := range slices {
                if err := service.FlattenBaseInfo(ctx, value.(int64)); err != nil {
                    return err
                }
            }
        }
    }
    return nil
}

func FlattenMerchant(ctx context.Context, input interface{}, cfg interface{}) error {
    if service, err := load(cfg); err != nil {
        return err
    } else {
        if slices, err := safeweb_lib_utils.MakeSliceFromInput(input); err != nil {
            return err
        } else {
            for _, value := range slices {
                if err := service.FlattenMerchant(ctx, value.(int64)); err != nil {
                    return err
                }
            }
        }
    }
    return nil
}

func FlattenTerminal(ctx context.Context, input interface{}, cfg interface{}) error {
    if service, err := load(cfg); err != nil {
        return err
    } else {
        if slices, err := safeweb_lib_utils.MakeSliceFromInput(input); err != nil {
            return err
        } else {
            for _, value := range slices {
                if err := service.FlattenTerminal(ctx, value.(int64)); err != nil {
                    return err
                }
            }
        }
    }
    return nil
}

func FlattenBank(ctx context.Context, input interface{}, cfg interface{}) error {
    if service, err := load(cfg); err != nil {
        return err
    } else {
        if slices, err := safeweb_lib_utils.MakeSliceFromInput(input); err != nil {
            return err
        } else {
            for _, value := range slices {
                if err := service.FlattenBank(ctx, value.(int64)); err != nil {
                    return err
                }
            }
        }
    }
    return nil
}

func FlattenBenefit(ctx context.Context, input interface{}, cfg interface{}) error {
    if service, err := load(cfg); err != nil {
        return err
    } else {
        if slices, err := safeweb_lib_utils.MakeSliceFromInput(input); err != nil {
            return err
        } else {
            for _, value := range slices {
                if err := service.FlattenBenefit(ctx, value.(int64)); err != nil {
                    return err
                }
            }
        }
    }
    return nil
}

func FlattenBlacklistCustomer(ctx context.Context, input interface{}, cfg interface{}) error {
    if service, err := load(cfg); err != nil {
        return err
    } else {
        if slices, err := safeweb_lib_utils.MakeSliceFromInput(input); err != nil {
            return err
        } else {
            for _, value := range slices {
                if err := service.FlattenBlacklistCustomer(ctx, value.(string)); err != nil {
                    return err
                }
            }
        }
    }
    return nil
}

func FlattenBlacklistTerminal(ctx context.Context, input interface{}, cfg interface{}) error {
    if service, err := load(cfg); err != nil {
        return err
    } else {
        if slices, err := safeweb_lib_utils.MakeSliceFromInput(input); err != nil {
            return err
        } else {
            for _, value := range slices {
                if err := service.FlattenBlacklistTerminal(ctx, value.(string)); err != nil {
                    return err
                }
            }
        }
    }
    return nil
}

func load(cfg interface{}) (*Service, error) {
    if conf, err := fetchConfig(cfg); err != nil {
        return nil, err
    } else {
        // Init connection
        dbConn, redisConn, err := config.InitConn(conf.App[config.AppAliasBackend], true, true)
        if err != nil {
            return nil, err
        } else {
            return InitializeService(dbConn, redisConn), nil
        }
    }
}

func fetchConfig(cfg interface{}) (conf *config.Config, err error) {
    switch cfg.(type) {
    case *config.Config:
        conf = cfg.(*config.Config)
        return
    case string:
        return config.Load()
    default:
        err = errors.New("Config Invalid Type")
        return
    }
}

func CallTask(ctx context.Context, input interface{}, cfg interface{}, f func(context.Context, interface{}, interface{}) error) error {
    return f(ctx, input, cfg)
}
