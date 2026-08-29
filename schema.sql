-- ============================================================================
-- QUANTITO ENGINE - MySQL Database Schema & Stored Procedures
-- Engine: InnoDB | Character Set: utf8mb4 | Collation: utf8mb4_unicode_ci
-- ============================================================================

CREATE DATABASE IF NOT EXISTS `quantito-engine`
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE `quantito-engine`;

SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS `asset_variants`;
DROP TABLE IF EXISTS `variants`;
DROP TABLE IF EXISTS `strategies`;
DROP TABLE IF EXISTS `experiment_assets`;
DROP TABLE IF EXISTS `experiments`;
SET FOREIGN_KEY_CHECKS = 1;

-- ============================================================================
-- 1. TABLE: EXPERIMENTS (Quantitative Experiments)
-- ============================================================================
CREATE TABLE `experiments` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `experiment` VARCHAR(100) NOT NULL,
  `description` VARCHAR(255) NULL,
  `params` JSON NOT NULL,
  `status` VARCHAR(20) NOT NULL DEFAULT 'ACTIVE', -- 'ACTIVE', 'PAUSED', 'COMPLETED'
  `notes` TEXT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_experiments_status` (`status`),
  INDEX `idx_experiments_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- 2. TABLE: EXPERIMENT_ASSETS (Asset Assignment to Experiment)
-- ============================================================================
CREATE TABLE `experiment_assets` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `experiment_id` INT UNSIGNED NOT NULL,
  `asset_id` INT UNSIGNED NOT NULL,
  `status` VARCHAR(20) NOT NULL DEFAULT 'ACTIVE', -- 'ACTIVE', 'EXCLUDED'
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_experiment_asset` (`experiment_id`, `asset_id`),
  INDEX `idx_experiment_assets_asset` (`asset_id`, `status`),
  INDEX `idx_experiment_assets_exp` (`experiment_id`, `status`),
  CONSTRAINT `fk_exp_assets_experiment`
    FOREIGN KEY (`experiment_id`) REFERENCES `experiments` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- 3. TABLE: STRATEGIES (Quantitative Strategies Catalog)
-- ============================================================================
CREATE TABLE `strategies` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL,
  `class` VARCHAR(100) NOT NULL, -- e.g.: 'fifty_day_avg', 'price_momentum'
  `params` JSON NOT NULL,
  `status` VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_strategies_class` (`class`),
  INDEX `idx_strategies_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- 4. TABLE: VARIANTS (Experiment Variants / Strategy Configurations)
-- ============================================================================
CREATE TABLE `variants` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `experiment_id` INT UNSIGNED NOT NULL,
  `strategy_id` INT UNSIGNED NOT NULL,
  `variant` VARCHAR(100) NOT NULL, -- e.g.: 'Variant A (SMA 50)', 'Variant B (RSI Momentum)'
  `notes` TEXT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_variants_experiment` (`experiment_id`),
  INDEX `idx_variants_strategy` (`strategy_id`),
  CONSTRAINT `fk_variants_experiment`
    FOREIGN KEY (`experiment_id`) REFERENCES `experiments` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_variants_strategy`
    FOREIGN KEY (`strategy_id`) REFERENCES `strategies` (`id`)
    ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- 5. TABLE: ASSET_VARIANTS (Signal Traceability and Historical Log by Asset)
-- ============================================================================
CREATE TABLE `asset_variants` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `variant_id` INT UNSIGNED NOT NULL,
  `experiment_id` INT UNSIGNED NOT NULL,
  `strategy_id` INT UNSIGNED NOT NULL,
  `asset_id` INT UNSIGNED NOT NULL,
  `signal` VARCHAR(10) NOT NULL, -- 'BUY', 'HOLD', 'SELL'
  `price` DECIMAL(18, 4) NOT NULL DEFAULT 0.0000,
  `confidence` DECIMAL(5, 4) NOT NULL DEFAULT 1.0000,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_asset_variants_asset` (`asset_id`, `created_at`),
  INDEX `idx_asset_variants_exp_asset` (`experiment_id`, `asset_id`, `created_at`),
  INDEX `idx_asset_variants_signal` (`signal`, `created_at`),
  CONSTRAINT `fk_av_variant`
    FOREIGN KEY (`variant_id`) REFERENCES `variants` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_av_experiment`
    FOREIGN KEY (`experiment_id`) REFERENCES `experiments` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_av_strategy`
    FOREIGN KEY (`strategy_id`) REFERENCES `strategies` (`id`)
    ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- STORED PROCEDURES
-- ============================================================================

DELIMITER $$

-- ----------------------------------------------------------------------------
-- 1. sp_find_experiments
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_find_experiments`$$
CREATE PROCEDURE `sp_find_experiments`(
  IN p_include_inactive BOOLEAN
)
BEGIN
  SELECT 
    `id`,
    `experiment`,
    `description`,
    `params`,
    `status`,
    `notes`
  FROM `experiments`
  WHERE p_include_inactive = TRUE OR `status` = 'ACTIVE'
  ORDER BY `id` ASC;
END$$

-- ----------------------------------------------------------------------------
-- 1.1 sp_get_experiment_by_id
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_get_experiment_by_id`$$
CREATE PROCEDURE `sp_get_experiment_by_id`(
  IN p_experiment_id INT UNSIGNED
)
BEGIN
  SELECT 
    `id`,
    `experiment`,
    `description`,
    `params`,
    `status`,
    `notes`
  FROM `experiments`
  WHERE `id` = p_experiment_id;
END$$

-- ----------------------------------------------------------------------------
-- 1.2 sp_get_experiment_by_name
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_get_experiment_by_name`$$
CREATE PROCEDURE `sp_get_experiment_by_name`(
  IN p_experiment VARCHAR(100)
)
BEGIN
  SELECT 
    `id`,
    `experiment`,
    `description`,
    `params`,
    `status`,
    `notes`
  FROM `experiments`
  WHERE `experiment` = p_experiment;
END$$

-- ----------------------------------------------------------------------------
-- 2. sp_find_experiment_assets
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_find_experiment_assets`$$
CREATE PROCEDURE `sp_find_experiment_assets`(
  IN p_experiment_id INT UNSIGNED
)
BEGIN
  SELECT 
    `id`,
    `experiment_id`,
    `asset_id`,
    `status`,
    `created_at`
  FROM `experiment_assets`
  WHERE `experiment_id` = p_experiment_id AND `status` = 'ACTIVE'
  ORDER BY `asset_id` ASC;
END$$

-- ----------------------------------------------------------------------------
-- 3. sp_find_experiments_by_asset
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_find_experiments_by_asset`$$
CREATE PROCEDURE `sp_find_experiments_by_asset`(
  IN p_asset_id INT UNSIGNED
)
BEGIN
  SELECT 
    e.`id`,
    e.`experiment`,
    e.`description`,
    e.`params`,
    e.`status`,
    e.`notes`
  FROM `experiments` e
  INNER JOIN `experiment_assets` ea ON ea.`experiment_id` = e.`id`
  WHERE ea.`asset_id` = p_asset_id
    AND ea.`status` = 'ACTIVE'
    AND e.`status` = 'ACTIVE'
  ORDER BY e.`id` ASC;
END$$

-- ----------------------------------------------------------------------------
-- 4. sp_find_variants_by_experiment
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_find_variants_by_experiment`$$
CREATE PROCEDURE `sp_find_variants_by_experiment`(
  IN p_experiment_id INT UNSIGNED
)
BEGIN
  SELECT 
    v.`id`,
    v.`experiment_id`,
    v.`strategy_id`,
    v.`variant`,
    v.`notes`,
    st.`name` AS `strategy_name`,
    st.`class`,
    st.`params`,
    st.`status` AS `strategy_status`
  FROM `variants` v
  INNER JOIN `strategies` st ON st.`id` = v.`strategy_id`
  WHERE (p_experiment_id = 0 OR v.`experiment_id` = p_experiment_id)
    AND st.`status` = 'ACTIVE'
  ORDER BY v.`id` ASC;
END$$

-- ----------------------------------------------------------------------------
-- 5. sp_save_asset_variant
-- Saves the strategy evaluation outcome (Signal, Price, Confidence).
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_save_asset_variant`$$
CREATE PROCEDURE `sp_save_asset_variant`(
  IN p_asset_id INT UNSIGNED,
  IN p_variant_id INT UNSIGNED,
  IN p_signal VARCHAR(10),
  IN p_price DECIMAL(18, 4),
  IN p_confidence DECIMAL(5, 4)
)
BEGIN
  DECLARE v_asset_variant_id INT UNSIGNED;
  DECLARE v_experiment_id INT UNSIGNED;
  DECLARE v_strategy_id INT UNSIGNED;

  SELECT `experiment_id`, `strategy_id`
  INTO v_experiment_id, v_strategy_id
  FROM `variants`
  WHERE `id` = p_variant_id;

  INSERT INTO `asset_variants` (
    `variant_id`,
    `experiment_id`,
    `strategy_id`,
    `asset_id`,
    `signal`,
    `price`,
    `confidence`,
    `created_at`
  ) VALUES (
    p_variant_id,
    v_experiment_id,
    v_strategy_id,
    p_asset_id,
    p_signal,
    p_price,
    p_confidence,
    NOW()
  );

  SET v_asset_variant_id = LAST_INSERT_ID();

  SELECT 
    `id`,
    `variant_id`,
    `experiment_id`,
    `strategy_id`,
    `asset_id`,
    `signal`,
    `price`,
    `confidence`,
    `created_at`
  FROM `asset_variants`
  WHERE `id` = v_asset_variant_id;
END$$

-- ----------------------------------------------------------------------------
-- 6. sp_find_asset_variants
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_find_asset_variants`$$
CREATE PROCEDURE `sp_find_asset_variants`(
  IN p_asset_id INT UNSIGNED,
  IN p_from_date DATETIME,
  IN p_to_date DATETIME
)
BEGIN
  SELECT 
    `id`,
    `variant_id`,
    `experiment_id`,
    `strategy_id`,
    `asset_id`,
    `signal`,
    `price`,
    `confidence`,
    `created_at`
  FROM `asset_variants`
  WHERE `asset_id` = p_asset_id
    AND (p_from_date IS NULL OR `created_at` >= p_from_date)
    AND (p_to_date IS NULL OR `created_at` <= p_to_date)
  ORDER BY `created_at` DESC;
END$$

DELIMITER ;
