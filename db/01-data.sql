-- ============================================================================
-- QUANTITO ENGINE - Initial Execution Setup Data
-- Database: quantito_engine
-- ============================================================================
USE `quantito_engine`;
-- Disable foreign key checks during batch insert
SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE `signals`;
TRUNCATE TABLE `execution_variants`;
TRUNCATE TABLE `executions`;
TRUNCATE TABLE `variants`;
TRUNCATE TABLE `strategies`;
TRUNCATE TABLE `experiments`;
TRUNCATE TABLE `assets`;
SET FOREIGN_KEY_CHECKS = 1;
-- ============================================================================
-- 0. ASSETS (Tradable Financial Instruments Catalog)
-- ============================================================================
INSERT INTO `assets` (
    `id`,
    `symbol`,
    `name`,
    `asset_type`,
    `exchange`,
    `currency`,
    `status`,
    `created_at`,
    `updated_at`
  )
VALUES (
    1,
    'AAPL',
    'Apple Inc.',
    'EQUITY',
    'NASDAQ',
    'USD',
    'ACTIVE',
    NOW(),
    NOW()
  ),
  (
    2,
    'MSFT',
    'Microsoft Corporation',
    'EQUITY',
    'NASDAQ',
    'USD',
    'ACTIVE',
    NOW(),
    NOW()
  ),
  (
    3,
    'NVDA',
    'NVIDIA Corporation',
    'EQUITY',
    'NASDAQ',
    'USD',
    'ACTIVE',
    NOW(),
    NOW()
  ),
  (
    4,
    'SPY',
    'SPDR S&P 500 ETF Trust',
    'ETF',
    'NYSE_ARCA',
    'USD',
    'ACTIVE',
    NOW(),
    NOW()
  ),
  (
    5,
    'QQQ',
    'Invesco QQQ Trust',
    'ETF',
    'NASDAQ',
    'USD',
    'ACTIVE',
    NOW(),
    NOW()
  ),
  (
    6,
    'BTC-USD',
    'Bitcoin USD',
    'CRYPTO',
    'CRYPTO',
    'USD',
    'ACTIVE',
    NOW(),
    NOW()
  ),
  (
    7,
    'ETH-USD',
    'Ethereum USD',
    'CRYPTO',
    'CRYPTO',
    'USD',
    'ACTIVE',
    NOW(),
    NOW()
  ),
  (
    8,
    'BTCUSDT',
    'Bitcoin Tether Spot',
    'CRYPTO',
    'BINANCE',
    'USDT',
    'ACTIVE',
    NOW(),
    NOW()
  );
-- ============================================================================
-- 1. STRATEGIES (Quantitative Strategies Catalog & Parameter Specs)
-- ============================================================================
INSERT INTO `strategies` (
    `id`,
    `name`,
    `class`,
    `param_definitions`,
    `status`,
    `created_at`,
    `updated_at`
  )
VALUES (
    1,
    '50-Day Moving Average',
    'fifty_day_avg',
    JSON_OBJECT(
      'fifty_day_avg',
      JSON_OBJECT(
        'type',
        'float',
        'required',
        true,
        'description',
        '50-day SMA baseline price'
      ),
      'current_price',
      JSON_OBJECT(
        'type',
        'float',
        'required',
        true,
        'description',
        'Current market price'
      ),
      'threshold',
      JSON_OBJECT(
        'type',
        'float',
        'default',
        0.02,
        'required',
        false,
        'description',
        'Minimum deviation band'
      )
    ),
    'ACTIVE',
    NOW(),
    NOW()
  ),
  (
    2,
    'Price Momentum',
    'price_momentum',
    JSON_OBJECT(
      'current_price',
      JSON_OBJECT(
        'type',
        'float',
        'required',
        true,
        'description',
        'Current market price'
      ),
      'pct_change',
      JSON_OBJECT(
        'type',
        'float',
        'required',
        true,
        'description',
        'Recent percentage change'
      ),
      'buy_threshold',
      JSON_OBJECT(
        'type',
        'float',
        'default',
        0.05,
        'required',
        false,
        'description',
        'Uptrend trigger'
      ),
      'sell_threshold',
      JSON_OBJECT(
        'type',
        'float',
        'default',
        -0.05,
        'required',
        false,
        'description',
        'Downtrend trigger'
      )
    ),
    'ACTIVE',
    NOW(),
    NOW()
  ),
  (
    3,
    'Clockwork HGB Regressor',
    'clockwork_hgb',
    JSON_OBJECT(
      'current_price',
      JSON_OBJECT(
        'type',
        'float',
        'required',
        true,
        'description',
        'Current bar open or close price'
      ),
      'gate_bps',
      JSON_OBJECT(
        'type',
        'float',
        'default',
        20.0,
        'required',
        true,
        'description',
        'q80 absolute score gate threshold in bps'
      ),
      'model_variant',
      JSON_OBJECT(
        'type',
        'string',
        'required',
        true,
        'description',
        'Frozen HGB model cut or label alignment variant identifier'
      ),
      'prediction_bps',
      JSON_OBJECT(
        'type',
        'float',
        'required',
        false,
        'description',
        'Precomputed or bar HGB prediction in bps'
      )
    ),
    'ACTIVE',
    NOW(),
    NOW()
  );
-- ============================================================================
-- 2. EXPERIMENTS (Quantitative Experiments Specification)
-- ============================================================================
INSERT INTO `experiments` (
    `id`,
    `name`,
    `description`,
    `status`,
    `execution_mode`,
    `start_date`,
    `end_date`,
    `frequency`,
    `max_capital`,
    `currency`,
    `params`,
    `notes`,
    `created_at`,
    `updated_at`
  )
VALUES (
    1,
    'Tech_Contraposition_SMA_vs_Momentum',
    'Estudio comparativo entre SMA de 50 dias y Price Momentum en acciones tecnologicas',
    'ACTIVE',
    'PAPER',
    '2026-01-01 00:00:00',
    NULL,
    '1d',
    25000.0000,
    'USD',
    JSON_OBJECT(
      'api_url',
      'http://localhost:8081',
      'benchmark',
      'SPY',
      'universe',
      'US_TECH_EQUITIES',
      'slippage_bps',
      5
    ),
    'Experimento de referencia para validacion y ejecucion del motor cuantitativo',
    NOW(),
    NOW()
  ),
  (
    2,
    'Crypto_Momentum_Baseline',
    'Estudio de momentum en criptomonedas (Pausado)',
    'PAUSED',
    'BACKTEST',
    '2026-01-01 00:00:00',
    '2026-06-30 23:59:59',
    '1h',
    10000.0000,
    'USD',
    JSON_OBJECT(
      'benchmark',
      'BTC-USD',
      'slippage_bps',
      10
    ),
    'Pausado para calibracion de umbrales',
    NOW(),
    NOW()
  ),
  (
    3,
    'Clockwork_BTC_2H_WalkForward',
    'Experimento Clockwork HGB de retorno logaritmico a 2h en BTCUSDT con contraposicion de cortes anuales de refit (Agosto, Octubre, Enero)',
    'ACTIVE',
    'BACKTEST',
    '2025-01-02 00:00:00',
    '2026-09-30 00:00:00',
    '2h',
    50000.0000,
    'USDT',
    JSON_OBJECT(
      'api_url',
      'http://localhost:8081',
      'benchmark',
      'BTCUSDT',
      'exchange',
      'Binance_Spot',
      'target_exposure',
      0.15,
      'hold_minutes',
      120,
      'execution_delay_minutes',
      1,
      'round_trip_cost_bps',
      5.0,
      'gate_quantile',
      0.80,
      'seed',
      432
    ),
    'POC BTC 2h (Clockwork): contraposicion de las 3 variantes de corte de entrenamiento y calibracion de gate q80',
    NOW(),
    NOW()
  ),
  (
    4,
    'Clockwork_Daily_CrossAsset_Adaptation',
    'Adaptacion diaria cross-asset (Stocks y ETFs) del modelo pooled Clockwork HGB comparando modos de alineacion de labels (execution_aligned vs source_label_emulation)',
    'ACTIVE',
    'BACKTEST',
    '2013-01-01 00:00:00',
    '2017-11-30 00:00:00',
    '1d',
    100000.0000,
    'USD',
    JSON_OBJECT(
      'api_url',
      'http://localhost:8081',
      'benchmark',
      'SPY',
      'universe',
      'US_STOCKS_AND_ETFS',
      'target_exposure',
      0.15,
      'hold_bars',
      4,
      'daily_delay_bars',
      1,
      'round_trip_cost_bps',
      5.0,
      'gate_quantile',
      0.80,
      'seed',
      432,
      'annualization_days',
      252
    ),
    'Adaptacion diaria cross-asset (Clockwork): contraposicion entre execution_aligned (t+1->t+5) y source_label_emulation (t->t+4)',
    NOW(),
    NOW()
  );
-- ============================================================================
-- 3. VARIANTS (Experiment Variants / Parametrized Strategy & Asset Setups)
-- ============================================================================
INSERT INTO `variants` (
  `id`, `experiment_id`, `strategy_id`, `asset_id`, `name`, `params`,
  `max_capital_pct`, `max_attempts`, `status`, `notes`, `created_at`, `updated_at`
) VALUES
-- Experiment 1: AAPL, MSFT, NVDA across SMA 50 (Strategy 1) and Price Momentum (Strategy 2)
(
  1,
  1,
  1,
  1,
  'AAPL - 50-Day SMA',
  JSON_OBJECT(
    'fifty_day_avg', 100.0,
    'current_price', 105.5,
    'threshold', 0.02
  ),
  (
    2,
    1,
    2,
    1,
    'AAPL - Price Momentum',
    JSON_OBJECT(
      'current_price',
      105.5,
      'pct_change',
      0.065,
      'buy_threshold',
      0.05,
      'sell_threshold',
      -0.05
    ),
    0.1667,
    3,
    'ACTIVE',
    'Estrategia de momentum basada en variacion porcentual reciente para AAPL',
    NOW(),
    NOW()
  ),
  (
    3,
    1,
    1,
    2,
    'MSFT - 50-Day SMA',
    JSON_OBJECT(
      'fifty_day_avg',
      410.0,
      'current_price',
      425.0,
      'threshold',
      0.02
    ),
    0.1667,
    3,
    'ACTIVE',
    'Estrategia tendencial basada en media movil de 50 dias para MSFT',
    NOW(),
    NOW()
  ),
  (
    4,
    1,
    2,
    2,
    'MSFT - Price Momentum',
    JSON_OBJECT(
      'current_price',
      425.0,
      'pct_change',
      0.055,
      'buy_threshold',
      0.05,
      'sell_threshold',
      -0.05
    ),
    0.1667,
    3,
    'ACTIVE',
    'Estrategia de momentum basada en variacion porcentual reciente para MSFT',
    NOW(),
    NOW()
  ),
  (
    5,
    1,
    1,
    3,
    'NVDA - 50-Day SMA',
    JSON_OBJECT(
      'fifty_day_avg',
      120.0,
      'current_price',
      132.0,
      'threshold',
      0.02
    ),
    0.1666,
    3,
    'ACTIVE',
    'Estrategia tendencial basada en media movil de 50 dias para NVDA',
    NOW(),
    NOW()
  ),
  (
    6,
    1,
    2,
    3,
    'NVDA - Price Momentum',
    JSON_OBJECT(
      'current_price',
      132.0,
      'pct_change',
      0.082,
      'buy_threshold',
      0.05,
      'sell_threshold',
      -0.05
    ),
    0.1666,
    3,
    'ACTIVE',
    'Estrategia de momentum basada en variacion porcentual reciente para NVDA',
    NOW(),
    NOW()
  ),
  -- Experiment 2: Crypto Momentum Baseline (Paused)
  (
    7,
    2,
    2,
    6,
    'BTC-USD - Crypto Conservative',
    JSON_OBJECT(
      'current_price',
      62000.0,
      'pct_change',
      0.02,
      'buy_threshold',
      0.08,
      'sell_threshold',
      -0.08
    ),
    1.0000,
    3,
    'ACTIVE',
    'Estrategia de momentum conservadora en cripto para BTC-USD',
    NOW(),
    NOW()
  ),
  -- Experiment 3: Clockwork BTC 2H WalkForward (BTCUSDT - Asset 8)
  (
    8,
    3,
    3,
    8,
    'BTCUSDT - Aug Fixed Postfit (aug_fixed_postfit)',
    JSON_OBJECT(
      'current_price',
      63500.0,
      'prediction_bps',
      26.5,
      'gate_bps',
      19.8,
      'model_variant',
      'aug_fixed_postfit'
    ),
    0.3333,
    3,
    'ACTIVE',
    'Entrena hasta el 3 de agosto de Y-1 y calibra el gate q80 post-fit entre el 4 de agosto y el 31 de diciembre de Y-1',
    NOW(),
    NOW()
  ),
  (
    9,
    3,
    3,
    8,
    'BTCUSDT - Oct Refit Postfit (oct_refit_postfit)',
    JSON_OBJECT(
      'current_price',
      63500.0,
      'prediction_bps',
      -24.0,
      'gate_bps',
      20.5,
      'model_variant',
      'oct_refit_postfit'
    ),
    0.3333,
    3,
    'ACTIVE',
    'Entrena hasta el 1 de octubre de Y-1 y calibra el gate q80 post-fit entre el 2 de octubre y el 31 de diciembre de Y-1',
    NOW(),
    NOW()
  ),
  (
    10,
    3,
    3,
    8,
    'BTCUSDT - Jan Refit Train90 (jan_refit_train90)',
    JSON_OBJECT(
      'current_price',
      63500.0,
      'prediction_bps',
      14.2,
      'gate_bps',
      21.0,
      'model_variant',
      'jan_refit_train90'
    ),
    0.3334,
    3,
    'ACTIVE',
    'Entrena hasta el 1 de enero de Y y calibra el gate q80 sobre los ultimos 90 dias calendario de entrenamiento',
    NOW(),
    NOW()
  ),
  -- Experiment 4: Clockwork Daily CrossAsset Adaptation
  (
    11,
    4,
    3,
    1,
    'AAPL - Execution Aligned (execution_aligned)',
    JSON_OBJECT(
      'current_price',
      175.0,
      'prediction_bps',
      25.4,
      'gate_bps',
      18.45,
      'model_variant',
      'execution_aligned'
    ),
    0.5000,
    3,
    'ACTIVE',
    'Modo principal diario: label y ejecucion alineados en t+1 -> t+5 (4 ruedas de holding con 1 rueda de delay)',
    NOW(),
    NOW()
  ),
  (
    12,
    4,
    3,
    1,
    'AAPL - Source Label Emulation (source_label_emulation)',
    JSON_OBJECT(
      'current_price',
      175.0,
      'prediction_bps',
      15.0,
      'gate_bps',
      18.45,
      'model_variant',
      'source_label_emulation'
    ),
    0.5000,
    3,
    'ACTIVE',
    'Modo comparacion diario: label entrenado en t -> t+4 pero ejecutado en t+1 -> t+5',
    NOW(),
    NOW()
  );