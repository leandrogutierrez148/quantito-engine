-- ============================================================================
-- QUANTITO ENGINE - Initial Execution Setup Data
-- Database: quantito-engine
-- ============================================================================

USE `quantito-engine`;

-- Disable foreign key checks during batch insert
SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE `signals`;
TRUNCATE TABLE `execution_variants`;
TRUNCATE TABLE `executions`;
TRUNCATE TABLE `variants`;
TRUNCATE TABLE `strategies`;
TRUNCATE TABLE `experiments`;
SET FOREIGN_KEY_CHECKS = 1;

-- ============================================================================
-- 1. STRATEGIES (Quantitative Strategies Catalog & Parameter Specs)
-- ============================================================================
INSERT INTO `strategies` (`id`, `name`, `class`, `param_definitions`, `status`, `created_at`, `updated_at`) VALUES
(
  1,
  '50-Day Moving Average',
  'fifty_day_avg',
  JSON_OBJECT(
    'fifty_day_avg', JSON_OBJECT('type', 'float', 'required', true, 'description', '50-day SMA baseline price'),
    'current_price', JSON_OBJECT('type', 'float', 'required', true, 'description', 'Current market price'),
    'threshold', JSON_OBJECT('type', 'float', 'default', 0.02, 'required', false, 'description', 'Minimum deviation band')
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
    'current_price', JSON_OBJECT('type', 'float', 'required', true, 'description', 'Current market price'),
    'pct_change', JSON_OBJECT('type', 'float', 'required', true, 'description', 'Recent percentage change'),
    'buy_threshold', JSON_OBJECT('type', 'float', 'default', 0.05, 'required', false, 'description', 'Uptrend trigger'),
    'sell_threshold', JSON_OBJECT('type', 'float', 'default', -0.05, 'required', false, 'description', 'Downtrend trigger')
  ),
  'ACTIVE',
  NOW(),
  NOW()
);

-- ============================================================================
-- 2. EXPERIMENTS (Quantitative Experiments Specification)
-- ============================================================================
INSERT INTO `experiments` (
  `id`, `name`, `description`, `status`, `execution_mode`,
  `start_date`, `end_date`, `frequency`, `max_capital`, `currency`,
  `assets`, `params`, `notes`, `created_at`, `updated_at`
) VALUES
(
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
  JSON_ARRAY('AAPL', 'MSFT', 'NVDA'),
  JSON_OBJECT(
    'benchmark', 'SPY',
    'universe', 'US_TECH_EQUITIES',
    'slippage_bps', 5
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
  JSON_ARRAY('BTC-USD', 'ETH-USD'),
  JSON_OBJECT(
    'benchmark', 'BTC-USD',
    'slippage_bps', 10
  ),
  'Pausado para calibracion de umbrales',
  NOW(),
  NOW()
);

-- ============================================================================
-- 3. VARIANTS (Experiment Variants / Parametrized Strategy Setups)
-- ============================================================================
INSERT INTO `variants` (
  `id`, `experiment_id`, `strategy_id`, `name`, `params`,
  `capital_allocation_pct`, `status`, `notes`, `created_at`, `updated_at`
) VALUES
(
  1,
  1,
  1,
  'Variant A (50-Day SMA)',
  JSON_OBJECT(
    'fifty_day_avg', 100.0,
    'current_price', 105.5,
    'threshold', 0.02
  ),
  0.5000,
  'ACTIVE',
  'Estrategia tendencial basada en media movil de 50 dias',
  NOW(),
  NOW()
),
(
  2,
  1,
  2,
  'Variant B (Price Momentum)',
  JSON_OBJECT(
    'current_price', 150.0,
    'pct_change', 0.065,
    'buy_threshold', 0.05,
    'sell_threshold', -0.05
  ),
  0.5000,
  'ACTIVE',
  'Estrategia de momentum basada en variacion porcentual reciente',
  NOW(),
  NOW()
),
(
  3,
  2,
  2,
  'Variant Crypto Conservative',
  JSON_OBJECT(
    'current_price', 150.0,
    'pct_change', 0.02,
    'buy_threshold', 0.08,
    'sell_threshold', -0.08
  ),
  1.0000,
  'ACTIVE',
  'Estrategia de momentum conservadora en cripto',
  NOW(),
  NOW()
);
